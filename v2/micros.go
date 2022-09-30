package micro

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/v2/pkg/config"
	"github.com/gidyon/micro/v2/pkg/conn"
	http_middleware "github.com/gidyon/micro/v2/pkg/middleware/http"
	redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"gorm.io/gorm"

	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Service contains API clients, connections and options for bootstrapping a micro-service.
type Service struct {
	cfg                      *config.Config
	logger                   grpclog.LoggerV2
	gormDBs                  map[string]*gorm.DB // uses gorm
	sqlDBs                   map[string]*sql.DB  // uses database/sql driver
	dbPoolOptions            map[string]*conn.DBConnPoolOptions
	redisClients             map[string]*redis.Client
	rediSearchClients        map[string]*redisearch.Client
	redisOptions             map[string]*redis.Options
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
	extServiceDialOptions    map[string][]grpc.DialOption
	unaryClientInterceptors  []grpc.UnaryClientInterceptor
	streamClientInterceptors []grpc.StreamClientInterceptor
	shutdowns                []func() error
	// timeouts
	httpServerReadTimeout  int
	httpServerWriteTimeout int
	initOnceFn             *sync.Once
	runOnceFn              *sync.Once
	nowFunc                func() time.Time
}

// NewService create a micro-service utility store by parsing data from config. Pass nil logger to use default logger
func NewService(_ context.Context, cfg *config.Config, grpcLogger grpclog.LoggerV2) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("nil config not allowed")
	}

	if cfg.StartupSleepSeconds() > 0 {
		time.Sleep(time.Duration(cfg.StartupSleepSeconds()) * time.Second)
	}

	var logger grpclog.LoggerV2

	if grpcLogger != nil {
		logger = grpcLogger
	} else {
		logger = NewLogger(cfg.ServiceName(), zerolog.TraceLevel)
	}

	svc := &Service{
		cfg:                      cfg,
		logger:                   logger,
		gormDBs:                  make(map[string]*gorm.DB),
		sqlDBs:                   make(map[string]*sql.DB),
		dbPoolOptions:            make(map[string]*conn.DBConnPoolOptions),
		redisClients:             make(map[string]*redis.Client),
		rediSearchClients:        make(map[string]*redisearch.Client),
		redisOptions:             make(map[string]*redis.Options),
		runtimeMuxEndpoint:       "",
		httpMiddlewares:          make([]http_middleware.Middleware, 0),
		httpMux:                  http.NewServeMux(),
		runtimeMux:               runtime.NewServeMux(),
		clientConn:               &grpc.ClientConn{},
		gRPCServer:               &grpc.Server{},
		externalServicesConn:     map[string]*grpc.ClientConn{},
		serveMuxOptions:          make([]runtime.ServeMuxOption, 0),
		serverOptions:            make([]grpc.ServerOption, 0),
		unaryInterceptors:        make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors:       make([]grpc.StreamServerInterceptor, 0),
		dialOptions:              make([]grpc.DialOption, 0),
		extServiceDialOptions:    map[string][]grpc.DialOption{},
		unaryClientInterceptors:  make([]grpc.UnaryClientInterceptor, 0),
		streamClientInterceptors: make([]grpc.StreamClientInterceptor, 0),
		shutdowns:                make([]func() error, 0),
		httpServerReadTimeout:    0,
		httpServerWriteTimeout:   0,
		initOnceFn:               &sync.Once{},
		runOnceFn:                &sync.Once{},
		nowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	return svc, nil
}

// Handler returns the http handler for the service
func (service *Service) Handler() http.Handler {
	return service.httpMux
}

// ServeMux returns the HTTP request multiplexer for the service
func (service *Service) ServeMux() *http.ServeMux {
	return service.httpMux
}

// SetServeMuxEndpoint sets the base pattern for the Http to gRPC handler
func (service *Service) SetServeMuxEndpoint(pattern string) {
	service.runtimeMuxEndpoint = pattern
}

// AddEndpoint registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (service *Service) AddEndpoint(pattern string, handler http.Handler) {
	if service.httpMux == nil {
		service.httpMux = http.NewServeMux()
	}
	service.httpMux.Handle(pattern, handler)
}

// AddEndpointFunc registers the handler function for the given pattern.
func (service *Service) AddEndpointFunc(pattern string, handleFunc http.HandlerFunc) {
	if service.httpMux == nil {
		service.httpMux = http.NewServeMux()
	}
	service.httpMux.HandleFunc(pattern, handleFunc)
}

// AddHTTPMiddlewares adds http middlewares to the service
func (service *Service) AddHTTPMiddlewares(middlewares ...http_middleware.Middleware) {
	service.httpMiddlewares = append(service.httpMiddlewares, middlewares...)
}

// AddExtServiceDialOptions add dials options to an external service service
func (service *Service) AddExtServiceDialOptions(externalService string, dialOptions ...grpc.DialOption) {
	service.extServiceDialOptions[strings.ToLower(externalService)] = append(service.extServiceDialOptions[strings.ToLower(externalService)], dialOptions...)
}

// AddGRPCDialOptions adds dial options to the service gRPC reverse proxy client
func (service *Service) AddGRPCDialOptions(dialOptions ...grpc.DialOption) {
	service.dialOptions = append(service.dialOptions, dialOptions...)
}

// AddGRPCServerOptions adds server options to the service gRPC server
func (service *Service) AddGRPCServerOptions(serverOptions ...grpc.ServerOption) {
	service.serverOptions = append(service.serverOptions, serverOptions...)
}

// AddGRPCStreamServerInterceptors adds stream interceptors to the service gRPC server
func (service *Service) AddGRPCStreamServerInterceptors(
	streamInterceptors ...grpc.StreamServerInterceptor,
) {
	service.streamInterceptors = append(
		service.streamInterceptors, streamInterceptors...,
	)
}

// AddGRPCUnaryServerInterceptors adds unary interceptors to the service gRPC server
func (service *Service) AddGRPCUnaryServerInterceptors(
	unaryInterceptors ...grpc.UnaryServerInterceptor,
) {
	service.unaryInterceptors = append(
		service.unaryInterceptors, unaryInterceptors...,
	)
}

// AddGRPCStreamClientInterceptors adds stream interceptors to the service gRPC reverse proxy client
func (service *Service) AddGRPCStreamClientInterceptors(
	streamInterceptors ...grpc.StreamClientInterceptor,
) {
	service.streamClientInterceptors = append(
		service.streamClientInterceptors, streamInterceptors...,
	)
}

// AddGRPCUnaryClientInterceptors adds unary interceptors to the service gRPC reverse proxy client
func (service *Service) AddGRPCUnaryClientInterceptors(
	unaryInterceptors ...grpc.UnaryClientInterceptor,
) {
	service.unaryClientInterceptors = append(
		service.unaryClientInterceptors, unaryInterceptors...,
	)
}

// AddRuntimeMuxOptions adds ServeMuxOption options to service gRPC reverse proxy client
// The options will be applied to the service runtime mux at startup
func (service *Service) AddRuntimeMuxOptions(serveMuxOptions ...runtime.ServeMuxOption) {
	if service.serveMuxOptions == nil {
		service.serveMuxOptions = make([]runtime.ServeMuxOption, 0)
	}
	service.serveMuxOptions = append(service.serveMuxOptions, serveMuxOptions...)
}

// GRPCDialOptions returns the service gRPC dial options
func (service *Service) GRPCDialOptions() []grpc.DialOption {
	return service.dialOptions
}

// Config returns the config used by the service
func (service *Service) Config() *config.Config {
	return service.cfg
}

// Logger returns grpc logger by the service
func (service *Service) Logger() grpclog.LoggerV2 {
	return service.logger
}

// RuntimeMux returns the HTTP request multiplexer for the service reverse proxy server
// gRPC services and methods are registered on this multiplxer.
// DO NOT register your anything on the returned muxer
// Use AddRuntimeMuxOptions method to register custom options
func (service *Service) RuntimeMux() *runtime.ServeMux {
	return service.runtimeMux
}

// ClientConn returns the underlying client connection to gRPC server used by reverse proxy
func (service *Service) ClientConn() *grpc.ClientConn {
	return service.clientConn
}

// GRPCServer returns the grpc server for the service
func (service *Service) GRPCServer() *grpc.Server {
	return service.gRPCServer
}

// GormDB returns the first gorm db with name "mysql"
func (service *Service) GormDB() *gorm.DB {
	return service.gormDBs["mysql"]
}

// GormDBByName returns the first gorm db with given name
func (service *Service) GormDBByName(name string) *gorm.DB {
	return service.gormDBs[name]
}

// GormDBs returns the underlying map for gorm dbs
func (service *Service) GormDBs() map[string]*gorm.DB {
	return service.gormDBs
}

// SQLDB returns the first database/sql db with name "mysql"
func (service *Service) SQLDB() *sql.DB {
	return service.sqlDBs["mysql"]
}

// SQLDBByName returns the first database/sql db with given name
func (service *Service) SQLDBByName(name string) *sql.DB {
	return service.sqlDBs[name]
}

// SQLDBs returns the underlying map for sql dbs
func (service *Service) SQLDBs() map[string]*sql.DB {
	return service.sqlDBs
}

// RedisClient returns the first redis client with name "redis"
func (service *Service) RedisClient() *redis.Client {
	return service.redisClients["redis"]
}

// RedisClientByName returns the first redis client with given name
func (service *Service) RedisClientByName(name string) *redis.Client {
	return service.redisClients[name]
}

// RediSearchClient returns the first redisearch client with name "redis"
func (service *Service) RediSearchClient() *redisearch.Client {
	return service.rediSearchClients["redis"]
}

// RediSearchClientByName returns the first redisearch client with given name
func (service *Service) RediSearchClientByName(name string) *redisearch.Client {
	return service.rediSearchClients[name]
}

// RedisClients returns the underlying map for redis client
func (service *Service) RedisClients() map[string]*redis.Client {
	return service.redisClients
}

// DialExternalService dials to an external service registered in config "externalServices" section using gRPC protocol
func (service *Service) DialExternalService(
	ctx context.Context, serviceName string, dialOptions ...grpc.DialOption,
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
		dopts = append(dopts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := conn.DialService(ctx, &conn.GRPCDialOptions{
		ServiceName: serviceInfo.Name(),
		Address:     serviceInfo.Address(),
		K8Service:   serviceInfo.K8Service(),
		DialOptions: dopts,
	})
	if err != nil {
		return nil, err
	}

	service.externalServicesConn[serviceInfo.Name()] = cc

	return cc, nil
}

// ExternalServiceConn returns the underlying rRPC connection to the external service found in config "externalServices"
// This method is prefered over DialExternalService at service startup
func (service *Service) ExternalServiceConn(serviceName string) (*grpc.ClientConn, error) {
	cc, ok := service.externalServicesConn[strings.ToLower(serviceName)]
	if !ok {
		return nil, errors.Errorf("no service exists with name: %s", serviceName)
	}
	return cc, nil
}

// SetDBConnPool sets connection pool options for sql database.
// If no database name is provided, the option will be applied to all sql databases
func (service *Service) SetDBConnPool(opt *conn.DBConnPoolOptions, dbNames ...string) {
	if service.dbPoolOptions == nil {
		service.dbPoolOptions = make(map[string]*conn.DBConnPoolOptions)
	}
	if len(dbNames) > 0 {
		for _, name := range dbNames {
			service.dbPoolOptions[name] = opt
		}
		return
	}
	for _, dbOptions := range service.Config().Databases() {
		if dbOptions.Type != config.SQLDBType {
			continue
		}
		service.dbPoolOptions[dbOptions.Metadata().Name()] = opt
	}
}

// SetRedisOptions sets redis options when starting redis client(s).
// If no client name is provided, the option will be applied to all redis client dbs
func (service *Service) SetRedisOptions(opt *redis.Options, dbNames ...string) {
	if service.redisOptions == nil {
		service.redisOptions = make(map[string]*redis.Options)
	}
	if len(dbNames) > 0 {
		for _, name := range dbNames {
			service.redisOptions[name] = opt
		}
		return
	}
	for _, dbOptions := range service.Config().Databases() {
		if dbOptions.Type != config.RedisDBType {
			continue
		}
		service.redisOptions[dbOptions.Metadata().Name()] = opt
	}
}

// SetHTTPServerReadTimout sets the read timeout to the service HTTP server
func (service *Service) SetHTTPServerReadTimout(sec int) {
	service.httpServerReadTimeout = sec
}

// SetHTTPServerWriteTimout sets the write timeout ftoor the service HTTP server
func (service *Service) SetHTTPServerWriteTimout(sec int) {
	service.httpServerWriteTimeout = sec
}

// SetNowFunc sets the function to be used when creating a new timestamp
func (service *Service) SetNowFunc(f func() time.Time) {
	service.nowFunc = f
}
