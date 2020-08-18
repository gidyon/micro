package micro

import (
	"context"
	"fmt"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gidyon/micro/pkg/conn"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (service *Service) openDBConn(ctx context.Context) error {
	var cfg = service.cfg

	if cfg.UseSQLDatabase() {
		sqlDBInfo := cfg.SQLDatabase()
		if sqlDBInfo.UseGorm() {
			gormDB, err := conn.OpenGormConn(&conn.DBOptions{
				Dialect:  sqlDBInfo.SQLDatabaseDialect(),
				Host:     sqlDBInfo.Host(),
				Port:     fmt.Sprintf("%d", sqlDBInfo.Port()),
				User:     sqlDBInfo.User(),
				Password: sqlDBInfo.Password(),
				Schema:   sqlDBInfo.Schema(),
				ConnPool: service.dbPoolOptions,
			})
			if err != nil {
				return err
			}

			service.gormDB = gormDB
		} else {
			sqlDB, err := conn.ToSQLDB(&conn.DBOptions{
				Dialect:  sqlDBInfo.SQLDatabaseDialect(),
				Host:     sqlDBInfo.Host(),
				Port:     fmt.Sprintf("%d", sqlDBInfo.Port()),
				User:     sqlDBInfo.User(),
				Password: sqlDBInfo.Password(),
				Schema:   sqlDBInfo.Schema(),
			})
			if err != nil {
				return err
			}

			service.sqlDB = sqlDB

			service.shutdowns = append(service.shutdowns, func() {
				sqlDB.Close()
			})
		}
	}

	return nil
}

func (service *Service) openRedisConn(ctx context.Context) error {
	var cfg = service.cfg

	if cfg.UseRedis() {
		redisDBInfo := cfg.RedisDatabase()

		service.redisClient = conn.NewRedisClient(&conn.RedisOptions{
			Address: redisDBInfo.Host(),
			Port:    fmt.Sprintf("%d", redisDBInfo.Port()),
		})

		service.shutdowns = append(service.shutdowns, func() {
			service.redisClient.Close()
		})

		if cfg.UseRediSearch() {
			service.rediSearchClient = redisearch.NewClient(
				redisDBInfo.Address(), cfg.ServiceName()+":index",
			)
		}
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
		if !srv.Available() {
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
