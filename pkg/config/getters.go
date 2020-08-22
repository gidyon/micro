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

// Databases returns list of all databases options
func (cfg *Config) Databases() []*DatabaseInfo {
	dbs := make([]*DatabaseInfo, 0, len(cfg.config.Databases))

	for _, db := range cfg.config.Databases {
		dbs = append(dbs, &DatabaseInfo{db})
	}

	return dbs
}

// DatabaseInfo contains parameters and connection information for a database
type DatabaseInfo struct {
	*databaseOptions
}

// Required is whether the database is required by the service
func (db *DatabaseInfo) Required() bool {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.Required
	}
	return false
}

// Address is the database address
func (db *DatabaseInfo) Address() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.Address
	}
	return ""
}

// User is the database user
func (db *DatabaseInfo) User() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.User
	}
	return ""
}

// Schema is the database schema
func (db *DatabaseInfo) Schema() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.Schema
	}
	return ""
}

// Password is the database password
func (db *DatabaseInfo) Password() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.Password
	}
	return ""
}

// UserFile is path to file containing database user
func (db *DatabaseInfo) UserFile() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.UserFile
	}
	return ""
}

// SchemaFile is path to file containing database schema
func (db *DatabaseInfo) SchemaFile() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.SchemaFile
	}
	return ""
}

// PasswordFile is the path to file containing database password
func (db *DatabaseInfo) PasswordFile() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.PasswordFile
	}
	return ""
}

// Metadata is read-only. Use getters to get individual metadata
func (db *DatabaseInfo) Metadata() {}

// UseGorm indicates whether the service will use Object Relational Mapper for database operations
func (db *DatabaseInfo) UseGorm() bool {
	if db != nil && db.databaseOptions != nil && db.databaseOptions.Metadata != nil {
		return db.databaseOptions.Metadata.Orm == "gorm"
	}
	return false
}

// UseRediSearch returns boolen that shows whether service uses redisearch
func (db *DatabaseInfo) UseRediSearch() bool {
	if db != nil && db.databaseOptions != nil && db.databaseOptions.Metadata != nil {
		return db.databaseOptions.Metadata.UseRediSearch
	}
	return false
}

// SQLDatabaseDialect returns the sql dialect for the sql database options
func (db *DatabaseInfo) SQLDatabaseDialect() string {
	if db != nil && db.databaseOptions != nil {
		return db.databaseOptions.Metadata.Dialect
	}
	return ""
}

// SQLDatabase returns the first sql database options for the service
func (cfg *Config) SQLDatabase() *DatabaseInfo {
	for _, db := range cfg.config.Databases {
		if db.Type == sqlDBType {
			return &DatabaseInfo{db}
		}
	}
	return nil
}

// SQLDatabaseByName returns the first sql database options with the given name
func (cfg *Config) SQLDatabaseByName(identifier string) *DatabaseInfo {
	for _, db := range cfg.config.Databases {
		if db.Type == sqlDBType && db.Metadata.Name == identifier {
			return &DatabaseInfo{db}
		}
	}
	return nil
}

// UseSQLDatabase indicates whether the service has sql database options
func (cfg *Config) UseSQLDatabase() bool {
	for _, db := range cfg.config.Databases {
		if db.Type == sqlDBType && db.Required {
			return true
		}
	}
	return false
}

// RedisDatabase returns the first redis database options for the service
func (cfg *Config) RedisDatabase() *DatabaseInfo {
	for _, db := range cfg.config.Databases {
		if db.Type == redisDBType {
			return &DatabaseInfo{db}
		}
	}
	return nil
}

// RedisDatabaseByName returns the first redis database options with the given name
func (cfg *Config) RedisDatabaseByName(name string) *DatabaseInfo {
	for _, db := range cfg.config.Databases {
		if db.Type == redisDBType && db.Metadata != nil {
			if db.Metadata.Name == name {
				return &DatabaseInfo{db}
			}
		}
	}
	return nil
}

// UseRedis returns whether service has redis options
func (cfg *Config) UseRedis() bool {
	for _, db := range cfg.config.Databases {
		if db.Type == redisDBType && db.Required {
			return true
		}
	}
	return false
}

// UseRediSearch returns whether service has redisearch options
func (cfg *Config) UseRediSearch() bool {
	for _, db := range cfg.config.Databases {
		if db.Type == redisDBType && db.Metadata.UseRediSearch {
			return true
		}
	}
	return false
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

// Required indicates whether the service is required
func (srv *ServiceInfo) Required() bool {
	return srv.externalServiceOptions.Required
}

// K8Service indicates whether the service is a k8s service
func (srv *ServiceInfo) K8Service() bool {
	return srv.externalServiceOptions.K8Service
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
