package micro

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gidyon/micro/v2/pkg/conn"
	http_middleware "github.com/gidyon/micro/v2/pkg/middleware/http"
	"github.com/gidyon/micro/v2/utils/tlsutil"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

// panic after encountering first non nil error
func handleErrs(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

// initializes service without starting it.
func (service *Service) init(ctx context.Context) {
	service.initOnceFn.Do(func() {
		handleErrs(
			service.openSQLDBConnections(ctx),
			service.openRedisConnections(ctx),
			service.openExternalConnections(ctx),
			service.initGRPC(ctx),
		)
	})
}

// Initialize initializes service without starting it.
func (service *Service) Initialize(ctx context.Context) {
	service.init(ctx)
}

// Start opens connection to databases and external services, afterwards starting grpc and http server to serve requests.
func (service *Service) Start(ctx context.Context, initFn func() error) {
	service.init(ctx)
	handleErrs(initFn(), service.run(ctx))
}

// starts the servers
func (service *Service) run(ctx context.Context) error {
	fn := func() error {
		defer func() {
			var err error
			for _, shutdown := range service.shutdowns {
				err = shutdown()
				if err != nil {
					service.logger.Errorln(err)
				}
			}
		}()

		// Apply optional middlewares
		if service.cfg.HttpOptions().CorsEnabled() {
			service.httpMiddlewares = append(service.httpMiddlewares, http_middleware.SupportCORS)
		}

		// Handles grpc gateway apis
		service.AddEndpoint(service.runtimeMuxEndpoint, service.runtimeMux)

		// Apply any middlewares to the handler
		handler := http_middleware.Apply(service.Handler(), service.httpMiddlewares...)

		var ghandler http.Handler

		// add grpc handler if TLS is enabled on service, will use same port
		if service.cfg.ServiceTLSEnabled() {
			ghandler = grpcHandlerFunc(service.GRPCServer(), handler)
		} else {
			ghandler = handler
		}

		httpServer := &http.Server{
			Addr:         fmt.Sprintf(":%d", service.cfg.ServicePort()),
			Handler:      ghandler,
			ReadTimeout:  time.Duration(service.httpServerReadTimeout) * time.Second,
			WriteTimeout: time.Duration(service.httpServerWriteTimeout) * time.Second,
		}

		// Graceful shutdown of server
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				service.logger.Warning("shutting down service ...")
				service.gRPCServer.Stop()
				log.Fatalln(httpServer.Shutdown(ctx))

				<-ctx.Done()
			}
		}()

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", service.cfg.ServicePort()))
		if err != nil {
			return errors.Wrap(err, "failed to create TCP listener for http server")
		}
		defer lis.Close()

		logMsgFn := func() {
			if !service.cfg.ServiceTLSEnabled() {
				service.logger.Infof(
					"<GRPC> running on port %d (insecure), <REST> server running on port %d (insecure)",
					service.cfg.GRPCPort(), service.cfg.ServicePort(),
				)
			} else {
				service.logger.Infof(
					"<gRPC> and <REST> server running on same port %d (secure)",
					service.cfg.ServicePort(),
				)
			}
		}

		logMsgFn()

		if !service.cfg.ServiceTLSEnabled() {
			glis, err := net.Listen("tcp", fmt.Sprintf(":%d", service.cfg.GRPCPort()))
			if err != nil {
				return errors.Wrap(err, "failed to create TCP listener for gRPC server")
			}
			defer glis.Close()

			// Note: The call to serve grpc must be inside a goroutine; don't do [go service.gRPCServer.Serve(glis)]
			go func() {
				err := service.gRPCServer.Serve(glis)
				if err != nil {
					service.logger.Errorln(err)
				}
			}()

			// Serve http insecurely
			return httpServer.Serve(lis)
		}

		cert, certPool, err := tlsutil.GetCert(service.Config().ServiceTLSCertFile(), service.Config().ServiceTLSKeyFile())
		if err != nil {
			return err
		}

		tlsConfig := &tls.Config{
			NextProtos:         []string{"h2", "http/1.1", "http/1.2"},
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS13,
			ClientAuth:         tls.VerifyClientCertIfGiven,
			ClientCAs:          certPool,
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: true,
		}

		// Serve tls
		return httpServer.Serve(tls.NewListener(lis, tlsConfig))
	}

	var err error
	service.runOnceFn.Do(func() {
		err = fn()
	})

	return err
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
	// ============================= Update runtime mux endpoint =============================
	if service.runtimeMuxEndpoint == "" {
		service.runtimeMuxEndpoint = "/"
	}

	// Apply servemux options to runtime muxer
	service.runtimeMux = runtime.NewServeMux(service.serveMuxOptions...)

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
		gPort = service.cfg.HTTPort()
	} else {
		service.dialOptions = append(service.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
		gPort = service.cfg.GRPCPort()
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
	unaryClientInterceptors = append(unaryClientInterceptors, service.unaryClientInterceptors...)

	// Add client streaming interceptos
	streamClientInterceptors := make([]grpc.StreamClientInterceptor, 0)
	streamClientInterceptors = append(streamClientInterceptors, service.streamClientInterceptors...)

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
