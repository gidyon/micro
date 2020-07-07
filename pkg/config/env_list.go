package config

const (
	// EnvServiceName iis the environment variable key for the service name
	EnvServiceName = "SERVICE_NAME"
	// EnvServicePort iis the environment variable key for the service port
	EnvServicePort = "SERVICE_PORT"
	// EnvRequireAuthentication iis the environment variable key for the whether service requires authentication
	EnvRequireAuthentication = "REQUIRE_AUTHENTICATION"
	// EnvAuthenticationAddress iis the environment variable key for the authentication address
	EnvAuthenticationAddress = "AUTHENTICATION_ADDRESS"
	// EnvAuthenticationHost iis the environment variable key for the authentication host
	EnvAuthenticationHost = "AUTHENTICATION_HOST"
	// EnvAuthenticationPort iis the environment variable key for the authentication port
	EnvAuthenticationPort = "AUTHENTICATION_PORT"
	// EnvAuthenticationCertFile iis the environment variable key for the authentication host tls cert file
	EnvAuthenticationCertFile = "AUTHENTICATION_TLS_CERT_FILE"
	// EnvAuthenticationServerName is the tls server name for the authentication service tls server name
	EnvAuthenticationServerName = "AUTHENTICATION_SERVER_NAME"
	// EnvAuthenticationCertURL iis the environment variable key for the authentication host tls cert url
	EnvAuthenticationCertURL = "AUTHENTICATION_TLS_CERT_URL"
	// EnvUseNotification iis the environment variable key for the notification host
	EnvUseNotification = "USE_NOTIFICATION"
	// EnvNotificationAddress iis the environment variable key for the notification host
	EnvNotificationAddress = "NOTIFICATION_ADDRESS"
	// EnvNotificationHost iis the environment variable key for the notification host
	EnvNotificationHost = "NOTIFICATION_HOST"
	// EnvNotificationPort iis the environment variable key for the notification port
	EnvNotificationPort = "NOTIFICATION_PORT"
	// EnvNotificationCertFile iis the environment variable key for the notification host tls cert file
	EnvNotificationCertFile = "NOTIFICATION_TLS_CERT_FILE"
	// EnvNotificationServerName is the tls server name for the notification service
	EnvNotificationServerName = "NOTIFICATION_SERVER_NAME"
	// EnvNotificationCertURL iis the environment variable key for the notification host tls cert url
	EnvNotificationCertURL = "NOTIFICATION_TLS_CERT_URL"
	// EnvLoggingLevel is the environment variable key for the logging level
	EnvLoggingLevel = "LOG_LEVEL"
	// EnvLoggingTimeFormat is the environment variable key for the logging time format
	EnvLoggingTimeFormat = "LOG_TIME_FORMAT"
	// EnvServiceTLSCertFile is the environment variable key for the  service tls cert file
	EnvServiceTLSCertFile = "TLS_CERT_FILE"
	// EnvServiceTLSKeyFile is the environment variable key for the service tls private key
	EnvServiceTLSKeyFile = "TLS_KEY_FILE"
	// EnvUseSQLDatabase is the environment variable key for whether service uses sql database
	EnvUseSQLDatabase = "USE_SQL_DATABASE"
	// EnvSQLDatabaseDialect is the environment variable key for the  sql database dialect
	EnvSQLDatabaseDialect = "SQL_DATABASE_DIALECT"
	// EnvSQLDatabaseORM is the environment variable key for the sql database orm
	EnvSQLDatabaseORM = "SQL_DATABASE_ORM"
	// EnvSQLDatabaseAddress is the environment variable key for the sql database address
	EnvSQLDatabaseAddress = "SQL_DATABASE_ADDRESS"
	// EnvSQLDatabaseHost is the environment variable key for the sql database host
	EnvSQLDatabaseHost = "SQL_DATABASE_HOST"
	// EnvSQLDatabasePort is the environment variable key for the sql database port
	EnvSQLDatabasePort = "SQL_DATABASE_PORT"
	// EnvSQLDatabaseUser is the environment variable key for the sql database user
	EnvSQLDatabaseUser = "SQL_DATABASE_USER"
	// EnvSQLDatabaseUserFile is the environment variable key for the sql database user file
	EnvSQLDatabaseUserFile = "SQL_DATABASE_USER_FILE"
	// EnvSQLDatabaseSchema is the environment variable key for the sql database schema
	EnvSQLDatabaseSchema = "SQL_DATABASE_DATABASE"
	// EnvSQLDatabaseSchemaFile is the environment variable key for the sql database schema file
	EnvSQLDatabaseSchemaFile = "SQL_DATABASE_DATABASE_FILE"
	// EnvSQLDatabasePassword is the environment variable key for the sql database password
	EnvSQLDatabasePassword = "SQL_DATABASE_PASSWORD"
	// EnvSQLDatabasePasswordFile is the environment variable key for the sql database password file
	EnvSQLDatabasePasswordFile = "SQL_DATABASE_PASSWORD_FILE"
	// EnvUseRedis is the environment variable key for the whether the service uses redis
	EnvUseRedis = "USE_REDIS"
	// EnvUseRediSearch is the environment variable key for whether the service uses redis search
	EnvUseRediSearch = "USE_REDISEARCH"
	// EnvRedisAddress is the environment variable key for the redis database address
	EnvRedisAddress = "REDIS_ADDRESS"
	// EnvRedisHost is the environment variable key for the redis database host
	EnvRedisHost = "REDIS_HOST"
	// EnvRedisPort is the environment variable key for the redis database port
	EnvRedisPort = "REDIS_PORT"
	// EnvRedisUser is the environment variable key for the redis database user
	EnvRedisUser = "REDIS_USER"
	// EnvRedisUserFile is the environment variable key for the redis database user file
	EnvRedisUserFile = "REDIS_USER_FILE"
	// EnvRedisSchema is the environment variable key for the redis database schema
	EnvRedisSchema = "REDIS_SCHEMA"
	// EnvRedisSchemaFile is the environment variable key for the redis database schema file
	EnvRedisSchemaFile = "REDIS_SCHEMA_FILE"
	// EnvRedisPassword is the environment variable key for the redis database password
	EnvRedisPassword = "REDIS_PASSWORD"
	// EnvRedisPasswordFile is the environment variable key for the redis database password file
	EnvRedisPasswordFile = "REDIS_PASSWORD_FILE"
)
