package micro

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/v2/pkg/config"
	"github.com/gidyon/micro/v2/pkg/conn"
	http_middleware "github.com/gidyon/micro/v2/pkg/middleware/http"
	redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
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
	jsonMarshalOptions       protojson.MarshalOptions
	jsonUnmarshalOptions     protojson.UnmarshalOptions
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
	shutdowns                []func()
	// timeouts
	httpServerReadTimeout  int
	httpServerWriteTimeout int
}

// NewService create a micro-service utility store by parsing data from config. Pass nil logger to use default logger
func NewService(ctx context.Context, cfg *config.Config, grpcLogger grpclog.LoggerV2) (*Service, error) {
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

	return &Service{
		cfg:                      cfg,
		logger:                   logger,
		gormDBs:                  make(map[string]*gorm.DB),
		sqlDBs:                   make(map[string]*sql.DB),
		dbPoolOptions:            make(map[string]*conn.DBConnPoolOptions),
		redisClients:             make(map[string]*redis.Client),
		rediSearchClients:        make(map[string]*redisearch.Client),
		redisOptions:             make(map[string]*redis.Options),
		httpMiddlewares:          make([]http_middleware.Middleware, 0),
		httpMux:                  http.NewServeMux(),
		runtimeMux:               runtime.NewServeMux(),
		serveMuxOptions:          make([]runtime.ServeMuxOption, 0),
		serverOptions:            make([]grpc.ServerOption, 0),
		unaryInterceptors:        make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors:       make([]grpc.StreamServerInterceptor, 0),
		unaryClientInterceptors:  make([]grpc.UnaryClientInterceptor, 0),
		streamClientInterceptors: make([]grpc.StreamClientInterceptor, 0),
		dialOptions:              make([]grpc.DialOption, 0),
		shutdowns:                make([]func(), 0),
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

// DialExternalService grpc dials to an external service
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
		dopts = append(dopts, grpc.WithInsecure())
	}

	return conn.DialService(ctx, &conn.GRPCDialOptions{
		ServiceName: serviceInfo.Name(),
		Address:     serviceInfo.Address(),
		K8Service:   serviceInfo.K8Service(),
		DialOptions: dopts,
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

// SetDBConnPool sets connection pool options for sql database.
// If no database name is provided, it will set for all sql databases
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
		if dbOptions.Type != config.SQLDBType && !dbOptions.Required() {
			continue
		}
		service.dbPoolOptions[dbOptions.Metadata().Name()] = opt
	}
}

// SetRedisOptions sets redis options when starting redis clients.
// If no client name is provided, it will set for all clients
func (service *Service) SetRedisOptions(opt *redis.Options, clientNames ...string) {
	if service.redisOptions == nil {
		service.redisOptions = make(map[string]*redis.Options)
	}
	if len(clientNames) > 0 {
		for _, name := range clientNames {
			service.redisOptions[name] = opt
		}
		return
	}
	for _, dbOptions := range service.Config().Databases() {
		if dbOptions.Type != config.RedisDBType && !dbOptions.Required() {
			continue
		}
		service.redisOptions[dbOptions.Metadata().Name()] = opt
	}
}

// SetHTTPServerReadTimout sets the read timeout for the http server
func (service *Service) SetHTTPServerReadTimout(sec int) {
	service.httpServerReadTimeout = sec
}

// SetHTTPServerWriteTimout sets the write timeout for the http server
func (service *Service) SetHTTPServerWriteTimout(sec int) {
	service.httpServerWriteTimeout = sec
}
