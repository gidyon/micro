package micro

import (
	"context"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/pkg/config"
	"github.com/gidyon/micro/pkg/conn"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (service *Service) openDBConn(ctx context.Context) error {
	var cfg = service.cfg

	for _, sqlDBInfo := range cfg.Databases() {
		if sqlDBInfo.Type != config.SQLDBType || !sqlDBInfo.Required() {
			continue
		}
		if sqlDBInfo.UseGorm() {
			gormDB, err := conn.OpenGormConn(&conn.DBOptions{
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

			service.gormDBs[sqlDBInfo.Metadata().Name()] = gormDB
		} else {
			sqlDB, err := conn.ToSQLDB(&conn.DBOptions{
				Dialect:  sqlDBInfo.SQLDatabaseDialect(),
				Address:  sqlDBInfo.Address(),
				User:     sqlDBInfo.User(),
				Password: sqlDBInfo.Password(),
				Schema:   sqlDBInfo.Schema(),
			})
			if err != nil {
				return err
			}

			service.sqlDBs[sqlDBInfo.Metadata().Name()] = sqlDB

			service.shutdowns = append(service.shutdowns, func() {
				sqlDB.Close()
			})
		}
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
			Network:  "tcp",
			Addr:     redisOptions.Address(),
			Password: redisOptions.Password(),
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
		if !srv.Required() {
			continue
		}

		dopts := make([]grpc.DialOption, 0)

		if !srv.Insecure() {
			creds, err := credentials.NewClientTLSFromFile(srv.TLSCertFile(), srv.ServerName())
			if err != nil {
				return errors.Wrapf(err, "failed to create tls config for %s service", srv.Name())
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
			return errors.Wrapf(err, "failed to create connection to service %s", srv.Name())
		}

		service.shutdowns = append(service.shutdowns, func() {
			externalServices[serviceName].Close()
		})
	}

	service.externalServicesConn = externalServices

	return nil
}
