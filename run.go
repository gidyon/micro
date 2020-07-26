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

// Start starts a grpc server and a http server that proxies requests to the grpc server.
func (service *Service) Start(ctx context.Context, initFn func() error) {
	// Bootstraps grpc server and client
	handleErr(service.initGRPC(ctx))
	// Execute init
	handleErr(initFn())
	// Start servers
	handleErr(service.run(ctx))
}

func (service *Service) run(ctx context.Context) error {
	defer service.shutdown()

	service.httpMux.Handle(service.runtimeMuxEndpoint, service.runtimeMux)

	handler := http_middleware.Apply(service.Handler(), service.httpMiddlewares...)

	var ghandler http.Handler

	if service.cfg.ServiceTLSEnabled() {
		ghandler = grpcHandlerFunc(service.GRPCServer(), handler)
	} else {
		ghandler = handler
	}

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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", service.cfg.ServicePort()))
	if err != nil {
		return errors.Wrap(err, "failed to create TCP listener for http server")
	}
	defer lis.Close()

	logMsgFn := func() {
		secureMsg := "secure"
		grpcPortMsg := ""
		if !service.cfg.ServiceTLSEnabled() {
			secureMsg = "insecure"
			grpcPortMsg = "8080(insecure-grpc) and"
		}
		service.logger.Infof(
			"<gRPC and REST> servers for service running on port %s %d(%s)",
			grpcPortMsg, service.cfg.ServicePort(), secureMsg,
		)
	}

	logMsgFn()

	if !service.cfg.ServiceTLSEnabled() {
		glis, err := net.Listen("tcp", ":8080")
		if err != nil {
			return errors.Wrap(err, "failed to create TCP listener for gRPC server")
		}
		defer glis.Close()

		// Serve grpc insecurely
		go service.gRPCServer.Serve(glis)

		// Serve http insecurely
		return httpServer.Serve(lis)
	}

	// Server http securely
	return httpServer.ServeTLS(lis, service.Config().ServiceTLSCertFile(), service.Config().ServiceTLSKeyFile())
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// initGRPC initialize gRPC server and client with registered client and server interceptors and options.
// The method must be called before registering anything on the gRPC server or passing options to the gRPC client.
// When this method has been called, subsequent calls to update interceptors becomes stale.
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
	var (
		gPort int
		err   error
	)
	if service.cfg.ServiceTLSEnabled() {
		creds, err := credentials.NewClientTLSFromFile(
			service.cfg.ServiceTLSCertFile(), service.cfg.ServiceTLSServerName())
		if err != nil {
			return errors.Wrapf(err,
				"failed to create tls config for %s service", service.cfg.ServiceTLSServerName())
		}
		service.dialOptions = append(service.dialOptions, grpc.WithTransportCredentials(creds))
		gPort = service.cfg.ServicePort()
	} else {
		service.dialOptions = append(service.dialOptions, grpc.WithInsecure())
		gPort = 8080 // for grpc
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
		Address:     fmt.Sprintf("localhost:%d", gPort),
		DialOptions: service.dialOptions,
		K8Service:   false,
	})
	if err != nil {
		return errors.Wrap(err, "client failed to dial to gRPC server")
	}

	// ============================= Initialize grpc server =============================
	// Add transport credentials if secure option is passed
	if service.cfg.ServiceTLSEnabled() {
		creds, err := credentials.NewServerTLSFromFile(
			service.cfg.ServiceTLSCertFile(), service.cfg.ServiceTLSKeyFile())
		if err != nil {
			return fmt.Errorf("failed to create grpc server tls credentials: %v", err)
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
