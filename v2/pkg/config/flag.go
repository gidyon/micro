package config

import (
	"flag"
)

func (cfg *config) setConfigFromFlag() {
	// Service section
	serviceName := flag.String("service-name", "", "Name of the service")
	httPort := flag.Int("http-port", 80, "Port number to bind the service http server")
	grpcPort := flag.Int("grpc-port", 8080, "Port number to bind to the service grpc server")
	startupSleepSeconds := flag.Int("startup-sleep-sec", 5, "Sleep period before starting the app/service")

	// SQL Database section
	useDB := flag.Bool("use-db", false, "Use RDBMS(mysql) database")
	dbORM := flag.String(
		"sqldb-orm", "gorm", "Object Relational Mapper (ORM) for querying database",
	)
	dbDialect := flag.String(
		"sqldb-dialect", "mysql", "SQL dialect to use",
	)
	dbName := flag.String(
		"sqldb-name", "mysql", "SQL Database name e.g mysql, postgres ...",
	)
	dbAddress := flag.String(
		"sqldb-address", "localhost:3306", "SQL Database address (can be ip address or domain name)",
	)
	dbUser := flag.String(
		"sqldb-user", "root", "SQL Database user",
	)
	dbUserFile := flag.String(
		"sqldb-user-file", "", "File location storing SQL Database user",
	)
	dbPassword := flag.String(
		"sqldb-password", "", "SQL Database password",
	)
	dbPasswordFile := flag.String(
		"sqldb-password-file", "", "File location storing SQL Database password",
	)
	dbSchema := flag.String(
		"sqldb-schema", "", "SQL Database schema to use",
	)
	dbSchemaFile := flag.String(
		"sqldb-schema-file", "", "File location storing SQL Database schema name",
	)

	// Redis section
	useRediSearch := flag.Bool(
		"use-redisearch", false, "Whether to use redis inverted search",
	)
	useRedis := flag.Bool(
		"use-redis", false, "Whether to use redis database",
	)
	redisAddress := flag.String(
		"redis-address", "localhost:6379", "Redis address (can be ip address of domain name)",
	)
	redisUser := flag.String(
		"redis-user", "", "Redis user",
	)
	redisUserFile := flag.String(
		"redis-user-file", "", "File location storing redis user name",
	)
	redisPassword := flag.String(
		"redis-password", "", "Redis password",
	)
	redisPasswordFile := flag.String(
		"redis-password-file", "", "File location storing redis password",
	)

	// Logging section
	logLevel := flag.Int(
		"log-level", -1, "Global log level",
	)

	// Service TLS certificate and key section
	tlsCertFile := flag.String(
		"tls-cert-file", "", "File location to TLS certificate for the service",
	)
	tlsKeyFile := flag.String(
		"tls-key-file", "", "File location to TLS private key for the service",
	)
	serverName := flag.String("tls-servername", "localhost", "Subject Alternative Name for tls certificate")
	insecure := flag.Bool("insecure", false, "Option for using insecure http")

	flag.Parse()

	cfgFromFlag := newConfig()

	// Service section
	cfgFromFlag.ServiceName = *serviceName
	cfgFromFlag.HTTPort = *httPort
	cfgFromFlag.GRPCPort = *grpcPort
	cfgFromFlag.StartupSleepSeconds = *startupSleepSeconds

	// logging section
	if *logLevel <= 5 && *logLevel >= -1 {
		cfgFromFlag.LogLevel = *logLevel
	}

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
