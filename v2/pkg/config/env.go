package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// setConfigFromEnv will set config from env variables but only if that config value is empty
func (cfg *config) setConfigFromEnv() error {
	cfgFromEnv := newConfig()

	// Service section
	if portStr := strings.TrimSpace(os.Getenv(EnvHTTPPort)); portStr != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(portStr, ":"))
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvHTTPPort)
		}
		cfgFromEnv.HTTPort = port
	}
	if portStr := strings.TrimSpace(os.Getenv(EnvGrpcPort)); portStr != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(portStr, ":"))
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvGrpcPort)
		}
		cfgFromEnv.GRPCPort = port
	}

	// Logging
	if logLevel := strings.TrimSpace(os.Getenv(EnvLogLevel)); logLevel != "" {
		logLevelInt64, err := strconv.Atoi(logLevel)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvLogLevel)
		}
		cfgFromEnv.LogLevel = int(logLevelInt64)
	}

	// Service TLS certificate and private key
	cfgFromEnv.Security.TLSCertFile = strings.TrimSpace(os.Getenv(EnvServiceTLSCertFile))
	cfgFromEnv.Security.TLSKeyFile = strings.TrimSpace(os.Getenv(EnvServiceTLSKeyFile))

	// SQLDatabase section
	if boolStr := strings.TrimSpace(os.Getenv(EnvUseSQLDatabase)); boolStr != "" {
		useDBBool, err := strconv.ParseBool(boolStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to boolean", EnvUseSQLDatabase)
		}
		cfgFromEnv.Database.SQLDatabase.Required = useDBBool
	}
	cfgFromEnv.Database.SQLDatabase.Address = strings.TrimSpace(os.Getenv(EnvSQLDatabaseAddress))
	cfgFromEnv.Database.SQLDatabase.User = strings.TrimSpace(os.Getenv(EnvSQLDatabaseUser))
	cfgFromEnv.Database.SQLDatabase.UserFile = strings.TrimSpace(os.Getenv(EnvSQLDatabaseUserFile))
	cfgFromEnv.Database.SQLDatabase.Password = strings.TrimSpace(os.Getenv(EnvSQLDatabasePassword))
	cfgFromEnv.Database.SQLDatabase.PasswordFile = strings.TrimSpace(os.Getenv(EnvSQLDatabasePasswordFile))
	cfgFromEnv.Database.SQLDatabase.Schema = strings.TrimSpace(os.Getenv(EnvSQLDatabaseSchema))
	cfgFromEnv.Database.SQLDatabase.SchemaFile = strings.TrimSpace(os.Getenv(EnvSQLDatabaseSchemaFile))
	cfgFromEnv.Database.SQLDatabase.Metadata.Orm = strings.TrimSpace(os.Getenv(EnvSQLDatabaseORM))
	cfgFromEnv.Database.SQLDatabase.Metadata.Dialect = strings.TrimSpace(os.Getenv(EnvSQLDatabaseDialect))

	// Redis section
	if userRediSearchStr := strings.TrimSpace(os.Getenv(EnvUseRediSearch)); userRediSearchStr != "" {
		useRediSearchBool, err := strconv.ParseBool(userRediSearchStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to boolean", EnvUseRediSearch)
		}
		cfgFromEnv.Database.RedisDatabase.Metadata.UseRediSearch = useRediSearchBool
	}
	if useRedisStr := strings.TrimSpace(os.Getenv(EnvUseRedis)); useRedisStr != "" {
		useRedisBool, err := strconv.ParseBool(useRedisStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to boolean", EnvUseRedis)
		}
		cfgFromEnv.Database.RedisDatabase.Required = useRedisBool
	}

	cfgFromEnv.Database.RedisDatabase.Address = strings.TrimSpace(os.Getenv(EnvRedisAddress))
	cfgFromEnv.Database.RedisDatabase.User = strings.TrimSpace(os.Getenv(EnvRedisUser))
	cfgFromEnv.Database.RedisDatabase.UserFile = strings.TrimSpace(os.Getenv(EnvRedisUserFile))
	cfgFromEnv.Database.RedisDatabase.Password = strings.TrimSpace(os.Getenv(EnvRedisPassword))
	cfgFromEnv.Database.RedisDatabase.PasswordFile = strings.TrimSpace(os.Getenv(EnvRedisPasswordFile))

	// External services section
	if cfg.ExternalServices == nil {
		cfg.ExternalServices = make([]*externalServiceOptions, 0)
	}

	// Update config with environmental config
	cfg.updateConfigWith(cfgFromEnv)

	return nil
}
