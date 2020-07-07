package micros

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/pkg/config"
	"github.com/gidyon/micro/pkg/conn"
	http_middleware "github.com/gidyon/micro/pkg/http"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"net/http"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Service contains API clients, connections and options for bootstrapping a micro-service.
type Service struct {
	cfg                      *config.Config
	logger                   grpclog.LoggerV2
	gormDB                   *gorm.DB // uses gorm
	sqlDB                    *sql.DB  // uses database/sql driver
	redisClient              *redis.Client
	rediSearchClient         *redisearch.Client
	runtimeMuxEndpoint       string
	httpMiddlewares          []http_middleware.Middleware
	httpMux                  *http.ServeMux
	runtimeMux               *runtime.ServeMux
	clientConn               *grpc.ClientConn
	gRPCServer               *grpc.Server
	externalServicesConn     map[string]*grpc.ClientConn
	serveMuxOptions          []runtime.ServeMuxOption
	serverOptions            []grpc.ServerOption
	unaryInterceptors        []grpc.UnaryServerInterceptor
	streamInterceptors       []grpc.StreamServerInterceptor
	dialOptions              []grpc.DialOption
	unaryClientInterceptors  []grpc.UnaryClientInterceptor
	streamClientInterceptors []grpc.StreamClientInterceptor
	shutdown                 func()
}

// NewService create a micro-service utility store by parsing data from config. Pass nil logger to use default logger
func NewService(ctx context.Context, cfg *config.Config, grpcLogger grpclog.LoggerV2) (*Service, error) {

	if cfg == nil {
		return nil, errors.New("nil config not allowed")
	}

	// Sleep if startup sleep is enabled
	if cfg.StartupSleepSeconds() > 0 {
		time.Sleep(time.Duration(cfg.StartupSleepSeconds()) * time.Second)
	}

	var (
		err              error
		gormDB           *gorm.DB
		sqlDB            *sql.DB
		redisClient      *redis.Client
		rediSearchClient *redisearch.Client
		externalServices = make(map[string]*grpc.ClientConn)
		logger           grpclog.LoggerV2
		cleanup          = make([]func(), 0)
	)

	if grpcLogger != nil {
		logger = grpcLogger
	} else {
		logger = NewLogger(cfg.ServiceName())
	}

	if cfg.UseSQLDatabase() {
		sqlDBInfo := cfg.SQLDatabase()
		if sqlDBInfo.UseGorm() {
			// Create a *sql.DB instance
			gormDB, err = conn.ToSQLDBUsingORM(&conn.DBOptions{
				Dialect:  sqlDBInfo.SQLDatabaseDialect(),
				Host:     sqlDBInfo.Host(),
				Port:     fmt.Sprintf("%d", sqlDBInfo.Port()),
				User:     sqlDBInfo.User(),
				Password: sqlDBInfo.Password(),
				Schema:   sqlDBInfo.Schema(),
			})
			if err != nil {
				return nil, err
			}
			sqlDB = nil

			cleanup = append(cleanup, func() {
				gormDB.Close()
			})

		} else {
			// Create a *sql.DB instance
			sqlDB, err = conn.ToSQLDB(&conn.DBOptions{
				Dialect:  sqlDBInfo.SQLDatabaseDialect(),
				Host:     sqlDBInfo.Host(),
				Port:     fmt.Sprintf("%d", sqlDBInfo.Port()),
				User:     sqlDBInfo.User(),
				Password: sqlDBInfo.Password(),
				Schema:   sqlDBInfo.Schema(),
			})
			if err != nil {
				return nil, err
			}
			gormDB = nil

			cleanup = append(cleanup, func() {
				sqlDB.Close()
			})
		}
	}

	if cfg.UseRedis() {
		redisDBInfo := cfg.RedisDatabase()

		// Creates a redis client
		redisClient = conn.NewRedisClient(&conn.RedisOptions{
			Address: redisDBInfo.Host(),
			Port:    fmt.Sprintf("%d", redisDBInfo.Port()),
		})

		cleanup = append(cleanup, func() {
			redisClient.Close()
		})

		if cfg.UseRediSearch() {
			// Create a redisearch client
			rediSearchClient = redisearch.NewClient(
				redisDBInfo.Address(), cfg.ServiceName()+":index",
			)
		}
	}

	// Remote services
	for _, srv := range cfg.ExternalServices() {
		if !srv.Available() {
			continue
		}

		dopts := make([]grpc.DialOption, 0)

		if !srv.Insecure() {
			creds, err := credentials.NewClientTLSFromFile(srv.TLSCertFile(), srv.ServerName())
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create tls config for %s service", srv.Name())
			}
			dopts = append(dopts, grpc.WithTransportCredentials(creds))
		} else {
			dopts = append(dopts, grpc.WithInsecure())
		}

		serviceName := strings.ToLower(srv.Name())
		externalServices[serviceName], err = conn.DialService(ctx, &conn.GRPCDialOptions{
			ServiceName: srv.Name(),
			Address:     srv.Address(),
			K8Service:   srv.K8Service(),
			DialOptions: dopts,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create connection to service %s", srv.Name())
		}

		cleanup = append(cleanup, func() {
			externalServices[serviceName].Close()
		})
	}

	return &Service{
		cfg:                      cfg,
		logger:                   logger,
		gormDB:                   gormDB,
		sqlDB:                    sqlDB,
		redisClient:              redisClient,
		rediSearchClient:         rediSearchClient,
		httpMiddlewares:          make([]http_middleware.Middleware, 0),
		httpMux:                  http.NewServeMux(),
		runtimeMux:               &runtime.ServeMux{},
		externalServicesConn:     externalServices,
		serveMuxOptions:          make([]runtime.ServeMuxOption, 0),
		serverOptions:            make([]grpc.ServerOption, 0),
		unaryInterceptors:        make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors:       make([]grpc.StreamServerInterceptor, 0),
		unaryClientInterceptors:  make([]grpc.UnaryClientInterceptor, 0),
		streamClientInterceptors: make([]grpc.StreamClientInterceptor, 0),
		dialOptions:              make([]grpc.DialOption, 0),
		shutdown: func() {
			for _, fn := range cleanup {
				fn()
			}
		},
	}, nil
}

// Handler returns the http handler for the service
func (service *Service) Handler() http.Handler {
	return service.httpMux
}

// ServeMux returns the internal ServeMux of the service
func (service *Service) ServeMux() *http.ServeMux {
	return service.httpMux
}

// SetServeMuxEndpoint sets the base path for servemux handler
func (service *Service) SetServeMuxEndpoint(path string) {
	service.runtimeMuxEndpoint = path
}

// AddEndpoint binds a handler to the service at provided path
func (service *Service) AddEndpoint(path string, handler http.Handler) {
	if service.httpMux == nil {
		service.httpMux = http.NewServeMux()
	}
	service.httpMux.Handle(path, handler)
}

// AddEndpointFunc works like http.HandleFunc
func (service *Service) AddEndpointFunc(path string, handleFunc http.HandlerFunc) {
	if service.httpMux == nil {
		service.httpMux = http.NewServeMux()
	}
	service.httpMux.HandleFunc(path, handleFunc)
}

// AddHTTPMiddlewares adds http middlewares to the service
func (service *Service) AddHTTPMiddlewares(middlewares ...http_middleware.Middleware) {
	service.httpMiddlewares = append(service.httpMiddlewares, middlewares...)
}

// AddGRPCDialOptions adds dial options to gRPC reverse proxy client
func (service *Service) AddGRPCDialOptions(dialOptions ...grpc.DialOption) {
	for _, dialOption := range dialOptions {
		service.dialOptions = append(service.dialOptions, dialOption)
	}
}

// AddGRPCServerOptions adds server options to gRPC server
func (service *Service) AddGRPCServerOptions(serverOptions ...grpc.ServerOption) {
	for _, serverOption := range serverOptions {
		service.serverOptions = append(service.serverOptions, serverOption)
	}
}

// AddGRPCStreamServerInterceptors adds stream interceptors to the gRPC server
func (service *Service) AddGRPCStreamServerInterceptors(
	streamInterceptors ...grpc.StreamServerInterceptor,
) {
	for _, streamInterceptor := range streamInterceptors {
		service.streamInterceptors = append(
			service.streamInterceptors, streamInterceptor,
		)
	}
}

// AddGRPCUnaryServerInterceptors adds unary interceptors to the gRPC server
func (service *Service) AddGRPCUnaryServerInterceptors(
	unaryInterceptors ...grpc.UnaryServerInterceptor,
) {
	for _, unaryInterceptor := range unaryInterceptors {
		service.unaryInterceptors = append(
			service.unaryInterceptors, unaryInterceptor,
		)
	}
}

// AddGRPCStreamClientInterceptors adds stream interceptors to the gRPC reverse proxy client
func (service *Service) AddGRPCStreamClientInterceptors(
	streamInterceptors ...grpc.StreamClientInterceptor,
) {
	for _, streamInterceptor := range streamInterceptors {
		service.streamClientInterceptors = append(
			service.streamClientInterceptors, streamInterceptor,
		)
	}
}

// AddGRPCUnaryClientInterceptors adds unary interceptors to the gRPC reverse proxy client
func (service *Service) AddGRPCUnaryClientInterceptors(
	unaryInterceptors ...grpc.UnaryClientInterceptor,
) {
	for _, unaryInterceptor := range unaryInterceptors {
		service.unaryClientInterceptors = append(
			service.unaryClientInterceptors, unaryInterceptor,
		)
	}
}

// AddServeMuxOptions adds servermux options to configure runtime mux
func (service *Service) AddServeMuxOptions(serveMuxOptions ...runtime.ServeMuxOption) {
	if service.serveMuxOptions == nil {
		service.serveMuxOptions = make([]runtime.ServeMuxOption, 0)
	}
	service.serveMuxOptions = append(service.serveMuxOptions, serveMuxOptions...)
}

// GRPCDialOptions returns the service grpc dial options
func (service *Service) GRPCDialOptions() []grpc.DialOption {
	return service.dialOptions
}

// Config returns the config for the service
func (service *Service) Config() *config.Config {
	return service.cfg
}

// Logger returns grpc logger for the service
func (service *Service) Logger() grpclog.LoggerV2 {
	return service.logger
}

// RuntimeMux returns the runtime muxer for the service
func (service *Service) RuntimeMux() *runtime.ServeMux {
	return service.runtimeMux
}

// ClientConn returns the underlying client connection to grpc server used by reverse proxy
func (service *Service) ClientConn() *grpc.ClientConn {
	return service.clientConn
}

// GRPCServer returns the grpc server
func (service *Service) GRPCServer() *grpc.Server {
	return service.gRPCServer
}

// GormDB returns a gorm db instance
func (service *Service) GormDB() *gorm.DB {
	return service.gormDB
}

// SQLDB returns database/sql db instance
func (service *Service) SQLDB() *sql.DB {
	return service.sqlDB
}

// RedisClient returns a redis client
func (service *Service) RedisClient() *redis.Client {
	return service.redisClient
}

// RediSearchClient returns redisearch client
func (service *Service) RediSearchClient() *redisearch.Client {
	return service.rediSearchClient
}

// DialExternalService dials to an external service
func (service *Service) DialExternalService(
	ctx context.Context, serviceName string, dialOptions []grpc.DialOption,
) (*grpc.ClientConn, error) {
	serviceInfo, err := service.Config().ExternalServiceByName(serviceName)
	if err != nil {
		return nil, err
	}

	dopts := append([]grpc.DialOption{}, dialOptions...)

	if !serviceInfo.Insecure() {
		creds, err := credentials.NewClientTLSFromFile(serviceInfo.TLSCertFile(), serviceInfo.ServerName())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create tls config for %s service", serviceInfo.Name())
		}
		dopts = append(dopts, grpc.WithTransportCredentials(creds))
	} else {
		dopts = append(dopts, grpc.WithInsecure())
	}

	return conn.DialService(ctx, &conn.GRPCDialOptions{
		ServiceName: serviceInfo.Name(),
		Address:     serviceInfo.Address(),
		K8Service:   serviceInfo.K8Service(),
		Insecure:    serviceInfo.Insecure(),
		DialOptions: dialOptions,
	})
}

// ExternalServiceConn returns the underlying grpc connection to the external service
func (service *Service) ExternalServiceConn(serviceName string) (*grpc.ClientConn, error) {
	cc, ok := service.externalServicesConn[strings.ToLower(serviceName)]
	if !ok {
		return nil, errors.Errorf("no service exists with name: %s", serviceName)
	}
	return cc, nil
}

// creates a http Muxer using runtime.NewServeMux
func newRuntimeMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				OrigName:     true,
				EmitDefaults: false,
			},
		),
	)
}
