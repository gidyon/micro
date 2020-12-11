package micro

import (
	"context"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/pkg/config"
	"github.com/gidyon/micro/pkg/conn"
	redis "github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (service *Service) openDBConn(ctx context.Context) error {
	var cfg = service.cfg

	for _, sqlDBInfo := range cfg.Databases() {
		if sqlDBInfo.Type != config.SQLDBType || !sqlDBInfo.Required() {
			continue
		}
		// open sql connection
		sqlDB, err := conn.ToSQLDB(&conn.DBOptions{
			Dialect:  sqlDBInfo.SQLDatabaseDialect(),
			Address:  sqlDBInfo.Address(),
			User:     sqlDBInfo.User(),
			Password: sqlDBInfo.Password(),
			Schema:   sqlDBInfo.Schema(),
			ConnPool: service.dbPoolOptions[sqlDBInfo.Metadata().Name()],
		})
		if err != nil {
			return err
		}

		service.sqlDBs[sqlDBInfo.Metadata().Name()] = sqlDB

		var gormDB *gorm.DB
		// open gorm connection
		switch strings.ToLower(sqlDBInfo.SQLDatabaseDialect()) {
		case "postgres":
			// gormDB, err = gorm.Open(postgres.New(postgres.Config{
			// 	Conn: sqlDB,
			// }), &gorm.Config{})
			// if err != nil {
			// 	return err
			// }
		default:
			// mysql connection
			gormDB, err = gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			if err != nil {
				return err
			}
		}

		service.gormDBs[sqlDBInfo.Metadata().Name()] = gormDB

		service.shutdowns = append(service.shutdowns, func() {
			sqlDB.Close()
		})
	}

	return nil
}

func (service *Service) openRedisConn(ctx context.Context) error {
	var cfg = service.cfg

	for _, redisOptions := range cfg.Databases() {
		if redisOptions.Type != config.RedisDBType && !redisOptions.Required() {
			continue
		}

		service.redisClients[redisOptions.Metadata().Name()] = conn.NewRedisClient(&redis.Options{
			Network:      "tcp",
			Addr:         redisOptions.Address(),
			Password:     redisOptions.Password(),
			MaxRetries:   5,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
			MinIdleConns: 10,
			MaxConnAge:   time.Hour,
		})

		if cfg.UseRediSearch() {
			service.rediSearchClients[redisOptions.Metadata().Name()] = redisearch.NewClient(
				redisOptions.Address(), cfg.ServiceName()+":index",
			)
		}

		service.shutdowns = append(service.shutdowns, func() {
			service.redisClients[redisOptions.Metadata().Name()].Close()
		})
	}

	return nil
}

func (service *Service) openExternalConn(ctx context.Context) error {
	var (
		err              error
		cfg              = service.cfg
		externalServices = make(map[string]*grpc.ClientConn)
	)

	// Remote services
	for _, srv := range cfg.ExternalServices() {
		srv := srv
		if !srv.Required() {
			continue
		}

		externalServices[srv.Name()], err = service.DialExternalService(ctx, srv.Name())
		if err != nil {
			return err
		}

		service.shutdowns = append(service.shutdowns, func() {
			externalServices[srv.Name()].Close()
		})
	}

	service.externalServicesConn = externalServices

	return nil
}
