package conn

import (
	"context"
	"database/sql"
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/balancer/roundrobin"

	"strings"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	// Imports mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DBOptions contains parameters for connecting to a SQL database
type DBOptions struct {
	Dialect  string
	Host     string
	Port     string
	User     string
	Password string
	Schema   string
}

// PortNumber return port with any colon(:) removed
func (opt *DBOptions) PortNumber() string {
	return strings.TrimPrefix(opt.Port, ":")
}

// ToSQLDBUsingORM opens a connection to a SQL database returning the gorm database client
func ToSQLDBUsingORM(opt *DBOptions) (*gorm.DB, error) {
	// add MySQL driver specific parameter to parse date/time
	param := "charset=utf8&parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		opt.User,
		opt.Password,
		opt.Host,
		opt.PortNumber(),
		opt.Schema,
		param,
	)

	dialect := func() string {
		if opt.Dialect == "" {
			return "mysql"
		}
		return opt.Dialect
	}()

	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "(gorm) failed to open connection to database")
	}

	return db, nil
}

// ToSQLDB opens a connection to an SQL database returning the database client
func ToSQLDB(opt *DBOptions) (*sql.DB, error) {
	// add MySQL driver specific parameter to parse date/time
	param := "charset=utf8&parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		opt.User,
		opt.Password,
		opt.Host,
		opt.PortNumber(),
		opt.Schema,
		param,
	)

	dialect := func() string {
		if opt.Dialect == "" {
			return "mysql"
		}
		return opt.Dialect
	}()

	sqlDB, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "(sql) failed to open connection to database")
	}

	return sqlDB, nil
}

// RedisOptions contains parameters for connecting to a redis database
type RedisOptions struct {
	Address string
	Port    string
}

// NewRedisClient opens a tcp connection to the redis database returning the client.
func NewRedisClient(opt *RedisOptions) *redis.Client {
	redisURL := func(a, b string) string {
		if a == "" {
			return b
		}
		return a
	}

	uri := fmt.Sprintf("%s:%s", opt.Address, strings.TrimPrefix(opt.Port, ":"))
	return redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    redisURL(uri, ":6379"),
	})
}

// GRPCDialOptions contains parameters for dialing a remote grpc service
type GRPCDialOptions struct {
	ServiceName string
	Address     string
	DialOptions []grpc.DialOption
	K8Service   bool
}

// DialService dials to a remote service returning the underlying grpc connection
func DialService(ctx context.Context, opt *GRPCDialOptions) (*grpc.ClientConn, error) {
	var (
		dopts = []grpc.DialOption{
			// Load balancer scheme
			grpc.WithBalancerName(roundrobin.Name),
			// Other interceptors
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(
					waitForReadyInterceptor,
				),
			),
		}
	)

	dopts = append(dopts, opt.DialOptions...)

	// Address for dialing the kubernetes service
	if opt.K8Service {
		opt.Address = strings.TrimSuffix(opt.Address, "dns:///")
		opt.Address = "dns:///" + opt.Address
	}

	return grpc.DialContext(ctx, opt.Address, dopts...)
}

func waitForReadyInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return invoker(ctx, method, req, reply, cc, append(opts, grpc.WaitForReady(true))...)
}
