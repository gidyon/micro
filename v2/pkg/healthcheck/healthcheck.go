package healthcheck

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/gidyon/micro/v2"
	"github.com/gidyon/micro/v2/pkg/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	// ProbeLiveNess indicates the health check probe is a liveness check i.e service is running correctly
	ProbeLiveNess = "liveness"
	// ProbeReadiness indicates the health check probe is a readiness check i.e service has started and can service requests
	ProbeReadiness = "readiness"
	// ProbeStartup indicates the health check probe is a startup check i.e service has started correctly
	ProbeStartup = "startup"
)

// ProbeOptions contains data and options required for doing healthcheck
type ProbeOptions struct {
	successMsg   string
	Service      *micro.Service
	AutoMigrator func() error
	Type         string
}

// RegisterProbe ...
func RegisterProbe(opt *ProbeOptions) http.HandlerFunc {
	if opt.AutoMigrator == nil {
		opt.AutoMigrator = func() error { return nil }
	}

	var (
		service = opt.Service
		cfg     = opt.Service.Config()
	)

	serviceNil := service == nil
	cfgNil := cfg == nil

	// apply defaults
	if !serviceNil && !cfgNil {
		switch opt.Type {
		case ProbeLiveNess:
			opt.successMsg = fmt.Sprintf("service %q is running correctly :)", cfg.ServiceName())
		case ProbeReadiness:
			opt.successMsg = fmt.Sprintf("service %q is ready :)", cfg.ServiceName())
		case ProbeStartup:
			opt.successMsg = fmt.Sprintf("service %q has started :)", cfg.ServiceName())
		default:
			opt.successMsg = fmt.Sprintf("service %q is ready and running :)", cfg.ServiceName())
		}
	}

	// Check only service or app internals and not external components
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle any panic
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("unexpected error: %v", err)
				fmt.Fprintln(w, errMsg)
			}
		}()

		var (
			mu   = &sync.Mutex{}
			errs = make([]string, 0)
			wg   = &sync.WaitGroup{}
			ctx  = r.Context()
			err  error
		)

		appendError := func(errMsg string) {
			mu.Lock()
			errs = append(errs, errMsg)
			mu.Unlock()
		}

		if serviceNil {
			w.WriteHeader(http.StatusExpectationFailed)
			_, err = w.Write([]byte("service is uninitialized"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if cfgNil {
			w.WriteHeader(http.StatusExpectationFailed)
			_, err = w.Write([]byte("service has no configuration options"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Check sql db connection
		if len(service.SQLDBs()) > 0 {
			for name, sqlDB := range service.SQLDBs() {
				name := name

				wg.Add(1)

				go func(sqlDB *sql.DB) {
					defer wg.Done()

					err = sqlDB.Ping()
					if err != nil {
						appendError(fmt.Sprintf("[%s] failed to ping sql database: %v", name, err))
						return
					}
				}(sqlDB)
			}
		}

		// Check gorm db connection
		if len(service.GormDBs()) > 0 {
			for name, gormDB := range service.GormDBs() {
				name := name

				wg.Add(1)

				go func(gormDB *gorm.DB) {
					defer wg.Done()

					sqlDB, err := gormDB.DB()
					if err != nil {
						appendError(fmt.Sprintf("[%s] failed to get sql database from gorm: %v", name, err))
						return
					}

					err = sqlDB.Ping()
					if err != nil {
						appendError(fmt.Sprintf("[%s] failed to ping sql database: %v", name, err))
						return
					}
				}(gormDB)
			}
		}

		// Check redis db connection
		if len(service.RedisClients()) > 0 {
			for name, redisClient := range service.RedisClients() {
				name := name

				wg.Add(1)

				go func(redisClient *redis.Client) {
					defer wg.Done()

					statusCMD := redisClient.Ping(ctx)
					if err := statusCMD.Err(); err != nil {
						appendError(fmt.Sprintf("[%s] failed to ping redis: %v", name, err))
						return
					}
				}(redisClient)
			}
		}

		// check external services
		for _, extSrv := range cfg.ExternalServices() {
			if !extSrv.Required() {
				continue
			}

			wg.Add(1)

			// dials concurrently
			go func(extSrv *config.ServiceInfo) {
				defer wg.Done()

				cc, err := service.DialExternalService(ctx, extSrv.Name())
				if err != nil {
					appendError(fmt.Sprintf("failed to connect to %s service: %v", extSrv.Name(), err))
				} else {
					defer cc.Close()
				}
			}(extSrv)
		}

		// wait for all dials and pings to complete
		wg.Wait()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)

		// Check errors from external components
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Fprintln(w, err)
			}
			return
		}

		fmt.Fprintln(w, opt.successMsg)
	}
}
