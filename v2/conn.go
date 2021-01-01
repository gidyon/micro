package micro

import (
	"context"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/v2/pkg/config"
	"github.com/gidyon/micro/v2/pkg/conn"
	redis "github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (service *Service) openSQLDBConnections(ctx context.Context) error {
	var cfg = service.cfg

	for _, sqlDBInfo := range cfg.Databases() {
		if sqlDBInfo.Type != config.SQLDBType || !sqlDBInfo.Required() {
			continue
		}

		clientName := sqlDBInfo.Metadata().Name()

		// open sql connection
		sqlDB, err := conn.OpenSQLDBConn(&conn.DBOptions{
			Dialect:  sqlDBInfo.SQLDatabaseDialect(),
			Address:  sqlDBInfo.Address(),
			User:     sqlDBInfo.User(),
			Password: sqlDBInfo.Password(),
			Schema:   sqlDBInfo.Schema(),
			ConnPool: service.dbPoolOptions[clientName],
		})
		if err != nil {
			return err
		}

		service.sqlDBs[clientName] = sqlDB

		var gormDB *gorm.DB
		// open gorm connection
		switch strings.ToLower(sqlDBInfo.SQLDatabaseDialect()) {
		case "postgres":
			gormDB, err = gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			if err != nil {
				return err
			}
		default:
			// mysql connection
			gormDB, err = gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			if err != nil {
				return err
			}
		}

		service.gormDBs[clientName] = gormDB

		service.shutdowns = append(service.shutdowns, func() {
			sqlDB.Close()
		})
	}

	return nil
}

func (service *Service) openRedisConnections(ctx context.Context) error {
	var cfg = service.cfg

	for _, redisOptions := range cfg.Databases() {
		if redisOptions.Type != config.RedisDBType || !redisOptions.Required() {
			continue
		}

		clientName := redisOptions.Metadata().Name()

		opts, ok := service.redisOptions[clientName]
		if !ok {
			// Default options
			opts = &redis.Options{
				Network:      "tcp",
				Addr:         redisOptions.Address(),
				Username:     redisOptions.User(),
				Password:     redisOptions.Password(),
				MaxRetries:   5,
				ReadTimeout:  time.Minute,
				WriteTimeout: time.Minute,
				MinIdleConns: 10,
				MaxConnAge:   time.Hour,
			}
		}

		opts.Addr = redisOptions.Address()
		opts.Username = redisOptions.User()
		opts.Password = redisOptions.Password()

		service.redisClients[clientName] = conn.OpenRedisConn(opts)

		if redisOptions.Metadata().UseRediSearch {
			service.rediSearchClients[clientName] = redisearch.NewClient(
				redisOptions.Address(), cfg.ServiceName()+":index",
			)
		}

		service.shutdowns = append(service.shutdowns, func() {
			service.redisClients[clientName].Close()
		})
	}

	return nil
}

func (service *Service) openExternalConnections(ctx context.Context) error {
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

		name := strings.ToLower(srv.Name())

		externalServices[name], err = service.DialExternalService(ctx, name)
		if err != nil {
			return err
		}

		service.shutdowns = append(service.shutdowns, func() {
			externalServices[name].Close()
		})
	}

	service.externalServicesConn = externalServices

	return nil
}
