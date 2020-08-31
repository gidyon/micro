package config

const (
	// EnvServiceName is environment variable key containing service name
	EnvServiceName = "SERVICE_NAME"
	// EnvHTTPPort is environment variable key containing service http port
	EnvHTTPPort = "HTTP_PORT"
	// EnvGrpcPort is environment variable key containing service grpc port
	EnvGrpcPort = "GRPC_PORT"
	// EnvLogLevel is environment variable key containing log level
	EnvLogLevel = "LOG_LEVEL"
	// EnvServiceTLSCertFile is environment variable key containing service tls cert file
	EnvServiceTLSCertFile = "TLS_CERT_FILE"
	// EnvServiceTLSKeyFile is environment variable key containing service tls private key
	EnvServiceTLSKeyFile = "TLS_KEY_FILE"
	// EnvUseSQLDatabase is environment variable key indicating whether service uses sql database
	EnvUseSQLDatabase = "USE_SQL_DATABASE"
	// EnvSQLDatabaseDialect is environment variable key containing sql database dialect
	EnvSQLDatabaseDialect = "SQL_DATABASE_DIALECT"
	// EnvSQLDatabaseORM is environment variable key containing sql database orm
	EnvSQLDatabaseORM = "SQL_DATABASE_ORM"
	// EnvSQLDatabaseAddress is environment variable key containing sql database address
	EnvSQLDatabaseAddress = "SQL_DATABASE_ADDRESS"
	// EnvSQLDatabaseUser is environment variable key containing sql database user
	EnvSQLDatabaseUser = "SQL_DATABASE_USER"
	// EnvSQLDatabaseUserFile is environment variable key containing sql database user file
	EnvSQLDatabaseUserFile = "SQL_DATABASE_USER_FILE"
	// EnvSQLDatabaseSchema is environment variable key containing sql database schema
	EnvSQLDatabaseSchema = "SQL_DATABASE_DATABASE"
	// EnvSQLDatabaseSchemaFile is environment variable key containing sql database schema file
	EnvSQLDatabaseSchemaFile = "SQL_DATABASE_DATABASE_FILE"
	// EnvSQLDatabasePassword is environment variable key containing sql database password
	EnvSQLDatabasePassword = "SQL_DATABASE_PASSWORD"
	// EnvSQLDatabasePasswordFile is environment variable key containing sql database password file
	EnvSQLDatabasePasswordFile = "SQL_DATABASE_PASSWORD_FILE"
	// EnvUseRedis is environment variable key indicating whether the service uses redis
	EnvUseRedis = "USE_REDIS"
	// EnvUseRediSearch is environment variable key indicating whether the service uses redis search
	EnvUseRediSearch = "USE_REDISEARCH"
	// EnvRedisAddress is environment variable key containing redis database address
	EnvRedisAddress = "REDIS_ADDRESS"
	// EnvRedisUser is environment variable key containing redis database user
	EnvRedisUser = "REDIS_USER"
	// EnvRedisUserFile is environment variable key containing redis database user file
	EnvRedisUserFile = "REDIS_USER_FILE"
	// EnvRedisSchema is environment variable key containing redis database schema
	EnvRedisSchema = "REDIS_SCHEMA"
	// EnvRedisSchemaFile is environment variable key containing redis database schema file
	EnvRedisSchemaFile = "REDIS_SCHEMA_FILE"
	// EnvRedisPassword is environment variable key containing redis database password
	EnvRedisPassword = "REDIS_PASSWORD"
	// EnvRedisPasswordFile is environment variable key containing redis database password file
	EnvRedisPasswordFile = "REDIS_PASSWORD_FILE"
)
