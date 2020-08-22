package config

import (
	"flag"
	"strconv"
	"strings"
)

func (cfg *config) setConfigFromFlag() {
	// Service section
	serviceName := flag.String("service-name", "", "Name of the service")
	servicePort := flag.String("service-port", "", "Port number to bind the service")
	startupSleepSeconds := flag.Int("startup-sleep-sec", 0, "Sleep period before starting the app/service")

	// RDBMS Database section
	useDB := flag.Bool("use-db", false, "Use RDBMS(mysql) database")
	dbORM := flag.String(
		"sqldb-orm", "gorm",
		"Object Relational Mapper (ORM) for querying database",
	)
	dbDialect := flag.String(
		"sqldb-dialect", "mysql",
		"SQL dialect to use",
	)
	dbName := flag.String(
		"sqldb-name", "",
		"RDBMS database name e.g mysql, postgres ...",
	)
	dbAddress := flag.String(
		"sqldb-address", "",
		"RDBMS database address (can be ip address or domain name)",
	)
	dbUser := flag.String(
		"sqldb-user", "",
		"RDBMS database user",
	)
	dbUserFile := flag.String(
		"sqldb-user-file", "",
		"File location storing RDBMS database user",
	)
	dbPassword := flag.String(
		"sqldb-password", "",
		"RDBMS database password",
	)
	dbPasswordFile := flag.String(
		"sqldb-password-file", "",
		"File location storing RDBMS database password",
	)
	dbSchema := flag.String(
		"sqldb-schema", "",
		"RDBMS database schema to use",
	)
	dbSchemaFile := flag.String(
		"sqldb-schema-file", "",
		"File location storing RDBMS database schema name",
	)

	// Redis section
	useRediSearch := flag.Bool(
		"use-redisearch", false,
		"Whether to use redis inverted search",
	)
	useRedis := flag.Bool(
		"use-redis", false,
		"Whether to use redis database",
	)
	redisAddress := flag.String(
		"redis-address", "",
		"Redis address (can be ip address of domain name)",
	)
	redisUser := flag.String(
		"redis-user", "",
		"Redis user",
	)
	redisUserFile := flag.String(
		"redis-user-file", "",
		"File location storing redis user name",
	)
	redisPassword := flag.String(
		"redis-password", "",
		"Redis password",
	)
	redisPasswordFile := flag.String(
		"redis-password-file", "",
		"File location storing redis password",
	)

	// Logging section
	logLevel := flag.Int(
		"log-level", 100,
		"Global log level",
	)
	logTimeFormat := flag.String(
		"log-time-format", "",
		"Time format for logger e.g 2006-01-02T15:04:05Z07:00",
	)

	// Service TLS certificate and key section
	tlsCertFile := flag.String(
		"tls-cert-file", "",
		"File location to TLS certificate for the service",
	)
	tlsKeyFile := flag.String(
		"tls-key-file", "",
		"File location to TLS private key for the service",
	)
	serverName := flag.String("servername", "", "Subject Alternative Name for tls certificate")
	insecure := flag.Bool("insecure", false, "Option for using insecure http")

	flag.Parse()

	cfgFromFlag := newConfig()

	// Service section
	cfgFromFlag.ServiceName = *serviceName
	if portStr := strings.TrimPrefix(*servicePort, ":"); portStr != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(portStr, ":"))
		if err != nil {
			panic("failed to parse service port")
		}
		cfgFromFlag.ServicePort = port
	}
	cfgFromFlag.StartupSleepSeconds = *startupSleepSeconds

	// logging section
	if *logLevel <= 5 && *logLevel >= -1 {
		cfgFromFlag.Logging.Level = *logLevel
	}
	cfgFromFlag.Logging.TimeFormat = *logTimeFormat

	// service security
	cfgFromFlag.Security.TLSCertFile = *tlsCertFile
	cfgFromFlag.Security.TLSKeyFile = *tlsKeyFile
	cfgFromFlag.Security.Insecure = *insecure
	cfgFromFlag.Security.ServerName = *serverName

	// RDMS section
	cfgFromFlag.Database.SQLDatabase.Metadata.Orm = *dbORM
	cfgFromFlag.Database.SQLDatabase.Metadata.Name = *dbName
	cfgFromFlag.Database.SQLDatabase.Metadata.Dialect = *dbDialect
	cfgFromFlag.Database.SQLDatabase.Required = *useDB
	cfgFromFlag.Database.SQLDatabase.Address = *dbAddress
	cfgFromFlag.Database.SQLDatabase.User = *dbUser
	cfgFromFlag.Database.SQLDatabase.UserFile = *dbUserFile
	cfgFromFlag.Database.SQLDatabase.Password = *dbPassword
	cfgFromFlag.Database.SQLDatabase.PasswordFile = *dbPasswordFile
	cfgFromFlag.Database.SQLDatabase.Schema = *dbSchema
	cfgFromFlag.Database.SQLDatabase.SchemaFile = *dbSchemaFile

	// Redis section
	cfgFromFlag.Database.RedisDatabase.Metadata.UseRediSearch = *useRediSearch
	cfgFromFlag.Database.RedisDatabase.Required = *useRedis
	cfgFromFlag.Database.RedisDatabase.Address = *redisAddress
	cfgFromFlag.Database.RedisDatabase.User = *redisUser
	cfgFromFlag.Database.RedisDatabase.UserFile = *redisUserFile
	cfgFromFlag.Database.RedisDatabase.Password = *redisPassword
	cfgFromFlag.Database.RedisDatabase.PasswordFile = *redisPasswordFile

	// Update config with config from flag
	cfg.updateConfigWith(cfgFromFlag)
}
