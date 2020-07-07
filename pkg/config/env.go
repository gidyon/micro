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
	cfgFromEnv.ServiceName = strings.TrimSpace(os.Getenv(EnvServiceName))
	if portStr := strings.TrimSpace(os.Getenv(EnvServicePort)); portStr != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(portStr, ":"))
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvServicePort)
		}
		cfgFromEnv.ServicePort = port
	}

	// Logging
	if logLevel := strings.TrimSpace(os.Getenv(EnvLoggingLevel)); logLevel != "" {
		logLevelInt64, err := strconv.Atoi(logLevel)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvLoggingLevel)
		}
		cfgFromEnv.Logging.Level = int(logLevelInt64)
	}
	cfgFromEnv.Logging.TimeFormat = strings.TrimSpace(os.Getenv(EnvLoggingTimeFormat))

	// Service TLS certificate and private key
	cfgFromEnv.Security.TLSCertFile = strings.TrimSpace(os.Getenv(EnvServiceTLSCertFile))
	cfgFromEnv.Security.TLSKeyFile = strings.TrimSpace(os.Getenv(EnvServiceTLSKeyFile))

	// SQLDatabase section
	if boolStr := strings.TrimSpace(os.Getenv(EnvUseSQLDatabase)); boolStr != "" {
		useDBBool, err := strconv.ParseBool(boolStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to boolean", EnvUseSQLDatabase)
		}
		cfgFromEnv.Databases.SQLDatabase.Required = useDBBool
	}
	cfgFromEnv.Databases.SQLDatabase.Address = strings.TrimSpace(os.Getenv(EnvSQLDatabaseAddress))
	cfgFromEnv.Databases.SQLDatabase.Host = strings.TrimSpace(os.Getenv(EnvSQLDatabaseHost))
	if portStr := strings.TrimSpace(os.Getenv(EnvSQLDatabasePort)); portStr != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(portStr, ":"))
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvSQLDatabasePort)
		}
		cfgFromEnv.Databases.SQLDatabase.Port = port
	}
	cfgFromEnv.Databases.SQLDatabase.User = strings.TrimSpace(os.Getenv(EnvSQLDatabaseUser))
	cfgFromEnv.Databases.SQLDatabase.UserFile = strings.TrimSpace(os.Getenv(EnvSQLDatabaseUserFile))
	cfgFromEnv.Databases.SQLDatabase.Password = strings.TrimSpace(os.Getenv(EnvSQLDatabasePassword))
	cfgFromEnv.Databases.SQLDatabase.PasswordFile = strings.TrimSpace(os.Getenv(EnvSQLDatabasePasswordFile))
	cfgFromEnv.Databases.SQLDatabase.Schema = strings.TrimSpace(os.Getenv(EnvSQLDatabaseSchema))
	cfgFromEnv.Databases.SQLDatabase.SchemaFile = strings.TrimSpace(os.Getenv(EnvSQLDatabaseSchemaFile))
	cfgFromEnv.Databases.SQLDatabase.Metadata.Orm = strings.TrimSpace(os.Getenv(EnvSQLDatabaseORM))
	cfgFromEnv.Databases.SQLDatabase.Metadata.Dialect = strings.TrimSpace(os.Getenv(EnvSQLDatabaseDialect))

	// Redis section
	if userRediSearchStr := strings.TrimSpace(os.Getenv(EnvUseRediSearch)); userRediSearchStr != "" {
		useRediSearchBool, err := strconv.ParseBool(userRediSearchStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to boolean", EnvUseRediSearch)
		}
		cfgFromEnv.Databases.RedisDatabase.Metadata.UseRediSearch = useRediSearchBool
	}
	if useRedisStr := strings.TrimSpace(os.Getenv(EnvUseRedis)); useRedisStr != "" {
		useRedisBool, err := strconv.ParseBool(useRedisStr)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to boolean", EnvUseRedis)
		}
		cfgFromEnv.Databases.RedisDatabase.Required = useRedisBool
	}

	cfgFromEnv.Databases.RedisDatabase.Address = strings.TrimSpace(os.Getenv(EnvRedisAddress))
	cfgFromEnv.Databases.RedisDatabase.Host = strings.TrimSpace(os.Getenv(EnvRedisHost))
	if portStr := strings.TrimSpace(os.Getenv(EnvRedisPort)); portStr != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(portStr, ":"))
		if err != nil {
			return errors.Wrapf(err, "failed to parse %s env to int", EnvRedisPort)
		}
		cfgFromEnv.Databases.RedisDatabase.Port = port
	}
	cfgFromEnv.Databases.RedisDatabase.User = strings.TrimSpace(os.Getenv(EnvRedisUser))
	cfgFromEnv.Databases.RedisDatabase.UserFile = strings.TrimSpace(os.Getenv(EnvRedisUserFile))
	cfgFromEnv.Databases.RedisDatabase.Password = strings.TrimSpace(os.Getenv(EnvRedisPassword))
	cfgFromEnv.Databases.RedisDatabase.PasswordFile = strings.TrimSpace(os.Getenv(EnvRedisPasswordFile))

	// External services section
	if cfg.ExternalServices == nil {
		cfg.ExternalServices = make([]*externalServiceOptions, 0)
	}

	// Update config with environmental config
	cfg.updateConfigWith(cfgFromEnv)

	return nil
}
