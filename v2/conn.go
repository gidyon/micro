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

		// Service pool options
		poolOptions, ok := service.dbPoolOptions[clientName]
		if !ok {
			poolOptions = &conn.DBConnPoolOptions{}
		}

		// Config pool options
		poolSettings := sqlDBInfo.PoolSettings()

		if poolOptions.MaxOpenConns == 0 {
			// Update pool option with one in config
			poolOptions.MaxOpenConns = poolSettings.MaxOpenConns()
		}
		if poolOptions.MaxIdleConns == 0 {
			// Update pool option with one in config
			poolOptions.MaxIdleConns = poolSettings.MaxIdleConns()
		}
		if poolOptions.MaxLifetime == 0 {
			// Update pool option with one in config
			poolOptions.MaxLifetime = time.Duration(time.Second * time.Duration(poolSettings.MaxConnLifetimeSeconds()))
		}

		// open sql connection
		sqlDB, err := conn.OpenSQLDBConn(&conn.DBOptions{
			Name:     sqlDBInfo.Metadata().Name(),
			Dialect:  sqlDBInfo.SQLDatabaseDialect(),
			Address:  sqlDBInfo.Address(),
			User:     sqlDBInfo.User(),
			Password: sqlDBInfo.Password(),
			Schema:   sqlDBInfo.Schema(),
			ConnPool: poolOptions,
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
			}), &gorm.Config{
				NowFunc: service.nowFunc,
			})
			if err != nil {
				return err
			}
		default:
			// mysql connection
			gormDB, err = gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
			}), &gorm.Config{
				NowFunc: service.nowFunc,
			})
			if err != nil {
				return err
			}
		}

		service.gormDBs[clientName] = gormDB

		service.shutdowns = append(service.shutdowns, func() error {
			return sqlDB.Close()
		})

		service.Logger().Infof("[CONNECTION TO SQL DATABASE MADE SUCCESSFULLY] [name: %s]", sqlDBInfo.Metadata().Name())
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

		service.shutdowns = append(service.shutdowns, func() error {
			return service.redisClients[clientName].Close()
		})

		service.Logger().Infof("[CONNECTION TO REDIS DATABASE MADE SUCCESSFULLY] [name: %s]", redisOptions.Metadata().Name())
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

		dopts := service.extServiceDialOptions[srv.Name()]
		if len(dopts) == 0 {
			dopts = make([]grpc.DialOption, 0)
		}

		externalServices[name], err = service.DialExternalService(ctx, name, dopts...)
		if err != nil {
			return err
		}

		service.shutdowns = append(service.shutdowns, func() error {
			return externalServices[name].Close()
		})

		service.Logger().Infof("[CONNECTION TO SERVICE MADE SUCCESSFULLY] [name: %s]", srv.Name())
	}

	service.externalServicesConn = externalServices

	return nil
}
