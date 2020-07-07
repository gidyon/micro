package micro

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/gidyon/micro/pkg/conn"
	http_middleware "github.com/gidyon/micro/pkg/http"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Start starts a http2 server that multiplexes gRPC and HTTP requests on the same port.
func (service *Service) Start(ctx context.Context, initFn func() error) {
	// Bootstrap service grpc
	handleErr(service.initGRPC(ctx))
	// Execute init
	handleErr(initFn())
	// Start http server for both grpc and RaESTful API
	handleErr(service.run(ctx))
}

func (service *Service) run(ctx context.Context) error {
	defer service.shutdown()

	service.httpMux.Handle(service.runtimeMuxEndpoint, service.runtimeMux)

	handler := http_middleware.Apply(service.Handler(), service.httpMiddlewares...)

	ghandler := grpcHandlerFunc(service.GRPCServer(), handler)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", service.cfg.ServicePort()),
		Handler:           ghandler,
		ReadTimeout:       time.Duration(5 * time.Second),
		ReadHeaderTimeout: time.Duration(5 * time.Second),
		WriteTimeout:      time.Duration(5 * time.Second),
	}

	// Graceful shutdown of server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			service.logger.Warning("shutting service...")
			httpServer.Shutdown(ctx)

			<-ctx.Done()
		}
	}()

	// Create TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", service.cfg.ServicePort()))
	if err != nil {
		return errors.Wrap(err, "failed to create TCP listener")
	}
	defer lis.Close()

	logMsgFn := func() {
		secureMsg := "secure"
		if !service.cfg.ServiceTLSEnabled() {
			secureMsg = "insecure"
		}
		service.logger.Infof(
			"<gRPC and REST> server for service running on port: %d (%s)", service.cfg.ServicePort(), secureMsg,
		)
	}

	logMsgFn()

	if !service.cfg.ServiceTLSEnabled() {
		return httpServer.Serve(lis)
	}

	return httpServer.ServeTLS(lis, service.Config().ServiceTLSCertFile(), service.Config().ServiceTLSKeyFile())
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// initGRPC initialize gRPC server and client with registered client and server interceptors and options.
// The method must be called before registering anything on the gRPC server or gRPC client connection.
// When this method has been called, subsequent calls to add interceptors and/or options will not update the service internals
func (service *Service) initGRPC(ctx context.Context) error {
	// ============================= Initialize runtime mux =============================
	if service.runtimeMuxEndpoint == "" {
		service.runtimeMuxEndpoint = "/"
	}

	// Apply servemux options to runtime muxer
	service.runtimeMux = runtime.NewServeMux(
		append(
			[]runtime.ServeMuxOption{},
			append(
				service.serveMuxOptions, runtime.WithMarshalerOption(
					runtime.MIMEWildcard,
					&runtime.JSONPb{
						OrigName:     true,
						EmitDefaults: true,
					},
				),
			)...,
		)...,
	)

	// ============================= Initialize grpc proxy client =============================
	var err error
	if service.cfg.ServiceTLSEnabled() {
		creds, err := credentials.NewClientTLSFromFile(service.cfg.ServiceTLSCertFile(), service.cfg.ServiceTLSServerName())
		if err != nil {
			return errors.Wrapf(err, "failed to create tls config for %s service", service.cfg.ServiceTLSServerName())
		}
		service.dialOptions = append(service.dialOptions, grpc.WithTransportCredentials(creds))
	} else {
		service.dialOptions = append(service.dialOptions, grpc.WithInsecure())
	}

	// Enable wait for ready RPCs
	waitForReadyUnaryInterceptor := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ctx, method, req, reply, cc, append(opts, grpc.WaitForReady(true))...)
	}

	// Add client unary interceptos
	unaryClientInterceptors := []grpc.UnaryClientInterceptor{waitForReadyUnaryInterceptor}
	for _, unaryInterceptor := range service.unaryClientInterceptors {
		unaryClientInterceptors = append(unaryClientInterceptors, unaryInterceptor)
	}

	// Add client streaming interceptos
	streamClientInterceptors := make([]grpc.StreamClientInterceptor, 0)
	for _, streamInterceptor := range service.streamClientInterceptors {
		streamClientInterceptors = append(streamClientInterceptors, streamInterceptor)
	}

	// Add inteceptors as dial option
	service.dialOptions = append(service.dialOptions, []grpc.DialOption{
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(unaryClientInterceptors...),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(streamClientInterceptors...),
		),
	}...)

	// client connection to the reverse gateway
	service.clientConn, err = conn.DialService(context.Background(), &conn.GRPCDialOptions{
		ServiceName: "self",
		Address:     fmt.Sprintf("localhost:%d", service.cfg.ServicePort()),
		DialOptions: service.dialOptions,
		K8Service:   false,
	})
	if err != nil {
		return errors.Wrap(err, "client failed to dial to gRPC server")
	}

	// ============================= Initialize grpc server =============================
	// Add transport credentials if secure
	if service.cfg.ServiceTLSEnabled() {
		creds, err := credentials.NewServerTLSFromFile(service.cfg.ServiceTLSCertFile(), service.cfg.ServiceTLSKeyFile())
		if err != nil {
			return fmt.Errorf("failed to create tls credentials: %v", err)
		}
		service.serverOptions = append(
			service.serverOptions, grpc.Creds(creds),
		)
	}

	// Append interceptors as server options
	service.serverOptions = append(
		service.serverOptions, grpc_middleware.WithUnaryServerChain(service.unaryInterceptors...))
	service.serverOptions = append(
		service.serverOptions, grpc_middleware.WithStreamServerChain(service.streamInterceptors...))

	service.gRPCServer = grpc.NewServer(service.serverOptions...)

	// register reflection on the gRPC server
	reflection.Register(service.gRPCServer)

	return nil
}
