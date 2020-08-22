package conn

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/balancer/roundrobin"

	"strings"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	// Imports mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DBConnPoolOptions contains options for customizing the connection pool
type DBConnPoolOptions struct {
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// DBOptions contains parameters for connecting to a SQL database
type DBOptions struct {
	Dialect  string
	Address  string
	User     string
	Password string
	Schema   string
	ConnPool *DBConnPoolOptions
}

// OpenGormConn open a gorm connection to the database
func OpenGormConn(opt *DBOptions) (*gorm.DB, error) {
	return ToSQLDBUsingORM(opt)
}

// ToSQLDBUsingORM opens a connection to a SQL database returning the gorm database client
func ToSQLDBUsingORM(opt *DBOptions) (*gorm.DB, error) {
	// Options should not be nil
	if opt == nil {
		return nil, errors.New("nil db options not allowed")
	}

	// add MySQL driver specific parameter to parse date/time
	param := "charset=utf8&parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		opt.User,
		opt.Password,
		opt.Address,
		opt.Schema,
		param,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, errors.Wrap(err, "(gorm) failed to open connection to mysql database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if opt.ConnPool != nil {
		if opt.ConnPool.MaxIdleConns != 0 {
			sqlDB.SetMaxIdleConns(opt.ConnPool.MaxIdleConns)
		}
		if opt.ConnPool.MaxOpenConns != 0 {
			sqlDB.SetMaxOpenConns(opt.ConnPool.MaxOpenConns)
		}
		if opt.ConnPool.MaxLifetime != 0 {
			sqlDB.SetConnMaxLifetime(opt.ConnPool.MaxLifetime)
		}
	}

	return db, nil
}

// ToSQLDB opens a connection to an SQL database returning the database client
func ToSQLDB(opt *DBOptions) (*sql.DB, error) {
	// Options should not be nil
	if opt == nil {
		return nil, errors.New("nil db options not allowed")
	}

	// add MySQL driver specific parameter to parse date/time
	param := "charset=utf8&parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		opt.User,
		opt.Password,
		opt.Address,
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

	if opt.ConnPool.MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(opt.ConnPool.MaxIdleConns)
	}
	if opt.ConnPool.MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(opt.ConnPool.MaxOpenConns)
	}
	if opt.ConnPool.MaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(opt.ConnPool.MaxLifetime)
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
