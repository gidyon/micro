package config

import (
	"strings"

	"github.com/pkg/errors"
)

// ServiceName returns the service name
func (cfg *Config) ServiceName() string {
	return cfg.config.ServiceName
}

// ServiceVersion returns the service version
func (cfg *Config) ServiceVersion() string {
	return cfg.config.ServiceVersion
}

// ServicePort returns the service port
func (cfg *Config) ServicePort() int {
	return cfg.config.ServicePort
}

// StartupSleepSeconds returns the startup sleep period
func (cfg *Config) StartupSleepSeconds() int {
	return cfg.config.StartupSleepSeconds
}

// ServiceTLSCertFile returns path to file containing tls certificate for the service
func (cfg *Config) ServiceTLSCertFile() string {
	return cfg.config.Security.TLSCertFile
}

// ServiceTLSKeyFile returns path to file containing tls private key for the service if any
func (cfg *Config) ServiceTLSKeyFile() string {
	return cfg.config.Security.TLSKeyFile
}

// ServiceTLSServerName returns tls server name of the service
func (cfg *Config) ServiceTLSServerName() string {
	return cfg.config.Security.ServerName
}

// ServiceTLSEnabled checks whether tls is enabled for the service
func (cfg *Config) ServiceTLSEnabled() bool {
	return !cfg.config.Security.Insecure
}

// Security prevent the struct field from being accidentally overriden
func (cfg *Config) Security() {}

// LogLevel returns log-level for logger
func (cfg *Config) LogLevel() int {
	return cfg.config.Logging.Level
}

// LogTimeFormat returns log time format for logger
func (cfg *Config) LogTimeFormat() string {
	return cfg.config.Logging.TimeFormat
}

// DisableLogger disables logging for the service
func (cfg *Config) DisableLogger() {
	cfg.config.Logging.Disabled = true
}

// Logging returns a boolean that show whether logging is enabled or disabled
func (cfg *Config) Logging() bool {
	return !cfg.config.Logging.Disabled
}

// Databases is only to make the embedded struct field read-only. Use Database getters instead
func (cfg *Config) Databases() {}

// DatabaseInfo contains parameters and connection information for a database
type DatabaseInfo struct {
	*databaseOptions
}

// Required is whether the database is required by the service
func (db *DatabaseInfo) Required() bool {
	return db != nil && db.databaseOptions.Required
}

// Address is the database address
func (db *DatabaseInfo) Address() string {
	return db.databaseOptions.Address
}

// Host is the database host
func (db *DatabaseInfo) Host() string {
	return db.databaseOptions.Host
}

// Port is the database port
func (db *DatabaseInfo) Port() int {
	return db.databaseOptions.Port
}

// User is the database user
func (db *DatabaseInfo) User() string {
	return db.databaseOptions.User
}

// Schema is the database schema
func (db *DatabaseInfo) Schema() string {
	return db.databaseOptions.Schema
}

// Password is the database password
func (db *DatabaseInfo) Password() string {
	return db.databaseOptions.Password
}

// UserFile is path to file containing database user
func (db *DatabaseInfo) UserFile() string {
	return db.databaseOptions.UserFile
}

// SchemaFile is path to file containing database schema
func (db *DatabaseInfo) SchemaFile() string {
	return db.databaseOptions.SchemaFile
}

// PasswordFile is the path to file containing database password
func (db *DatabaseInfo) PasswordFile() string {
	return db.databaseOptions.PasswordFile
}

// Metadata is read-only. Use getters to get individual metadata
func (db *DatabaseInfo) Metadata() {}

// UseGorm indicates whether the service will use Object Relational Mapper for database operations
func (db *DatabaseInfo) UseGorm() bool {
	return db.databaseOptions.Metadata.Orm == "gorm"
}

// UseRediSearch returns boolen that shows whether to use redisearch
func (db *DatabaseInfo) UseRediSearch() bool {
	return db.databaseOptions.Metadata.UseRediSearch
}

// SQLDatabaseDialect returns the dialect to used for interacting with sql database
func (db *DatabaseInfo) SQLDatabaseDialect() string {
	return db.databaseOptions.Metadata.Dialect
}

// SQLDatabase is the service sql database
func (cfg *Config) SQLDatabase() *DatabaseInfo {
	return &DatabaseInfo{cfg.config.Databases.SQLDatabase}
}

// UseSQLDatabase indicates whether the service will use mysql database
func (cfg *Config) UseSQLDatabase() bool {
	return cfg.config.Databases.SQLDatabase.Required
}

// RedisDatabase is the service redis database
func (cfg *Config) RedisDatabase() *DatabaseInfo {
	return &DatabaseInfo{cfg.config.Databases.RedisDatabase}
}

// UseRedis returns boolean that shows whether to use redis
func (cfg *Config) UseRedis() bool {
	return cfg.config.Databases.RedisDatabase.Required
}

// UseRediSearch returns boolean that shows whether to use redisearch
func (cfg *Config) UseRediSearch() bool {
	return cfg.config.Databases.RedisDatabase.Metadata.UseRediSearch
}

// ServiceInfo contains metadata and discovery information for a service
type ServiceInfo struct {
	*externalServiceOptions
}

// ExternalServices returns the list of available services
func (cfg *Config) ExternalServices() []*ServiceInfo {
	srvsInfo := make([]*ServiceInfo, 0, len(cfg.config.ExternalServices))
	for _, extSrv := range cfg.config.ExternalServices {
		srvsInfo = append(srvsInfo, &ServiceInfo{extSrv})
	}
	return srvsInfo
}

// Available indicates whether the service is required/available
func (srv *ServiceInfo) Available() bool {
	return srv.externalServiceOptions.Available
}

// K8Service indicates whether the service is a k8s service
func (srv *ServiceInfo) K8Service() bool {
	return srv.externalServiceOptions.K8Service
}

// Type returns the type of the service
func (srv *ServiceInfo) Type() string {
	return srv.externalServiceOptions.Type
}

// Name returns the name of the service
func (srv *ServiceInfo) Name() string {
	return srv.externalServiceOptions.Name
}

// Address returns the network address of the service
func (srv *ServiceInfo) Address() string {
	return srv.externalServiceOptions.Address
}

// Host returns the host name of the service
func (srv *ServiceInfo) Host() string {
	return srv.externalServiceOptions.Host
}

// Port returns the host name of the service
func (srv *ServiceInfo) Port() int {
	return srv.externalServiceOptions.Port
}

// TLSCertFile returns the service tls certificate file path
func (srv *ServiceInfo) TLSCertFile() string {
	return srv.externalServiceOptions.TLSCertFile
}

// ServerName returns the server name registered in tls certificare for service
func (srv *ServiceInfo) ServerName() string {
	return srv.externalServiceOptions.ServerName
}

// Insecure returns whether the service can be dialed over insecure http
func (srv *ServiceInfo) Insecure() bool {
	return srv.externalServiceOptions.Insecure
}

// ExternalServiceByName first service whose name matches the passed service name.
// The name comparison is case-insentive.
func (cfg *Config) ExternalServiceByName(serviceName string) (*ServiceInfo, error) {
	for _, srv := range cfg.ExternalServices() {
		if strings.ToLower(srv.Name()) == strings.ToLower(serviceName) {
			return srv, nil
		}
	}
	return nil, errors.Errorf("no service found with name: %s", serviceName)
}

// ExternalServiceByType first service whose type matches the passed service type.
// The type comparison is case-insentive.
func (cfg *Config) ExternalServiceByType(serviceType string) (*ServiceInfo, error) {
	for _, srv := range cfg.ExternalServices() {
		if strings.ToLower(srv.Type()) == strings.ToLower(serviceType) {
			return srv, nil
		}
	}
	return nil, errors.Errorf("no service found with type: %s", serviceType)
}
