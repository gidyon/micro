// Package config contains options for bootstrapping a service dependencies
package config

import (
	"github.com/pkg/errors"
)

type securityOptions struct {
	TLSCertFile string `yaml:"tlsCert"`
	TLSKeyFile  string `yaml:"tlsKey"`
	ServerName  string `yaml:"serverName"`
	Insecure    bool   `yaml:"insecure"`
}

type poolSettings struct {
	MaxOpenConns           uint `yaml:"maxOpenConns"`
	MaxIdleConns           uint `yaml:"maxIdleConns"`
	MaxConnLifetimeSeconds uint `yaml:"maxConnLifetimeSeconds"`
}

type dbMetadata struct {
	Name          string `yaml:"name"`
	Dialect       string `yaml:"dialect"`
	Orm           string `yaml:"orm"`
	UseRediSearch bool   `yaml:"useRediSearch"`
}

// databaseOptions contains parameters that open connection to a database
type databaseOptions struct {
	Required     bool          `yaml:"required"`
	Type         string        `yaml:"type"`
	Address      string        `yaml:"address"`
	User         string        `yaml:"user"`
	Schema       string        `yaml:"schema"`
	Password     string        `yaml:"password"`
	UserFile     string        `yaml:"userFile"`
	SchemaFile   string        `yaml:"schemaFile"`
	PasswordFile string        `yaml:"passwordFile"`
	PoolSettings *poolSettings `yaml:"poolSettings"`
	Metadata     *dbMetadata   `yaml:"metadata"`
}

// externalServiceOptions contains information to connect to a remote service
type externalServiceOptions struct {
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
	K8Service   bool   `yaml:"k8service"`
	Address     string `yaml:"address"`
	TLSCertFile string `yaml:"tlsCert"`
	ServerName  string `yaml:"serverName"`
	Insecure    bool   `yaml:"insecure"`
}

type httpOptions struct {
	CorsEnabled bool `yaml:"corsEnabled"`
}

// config contains configuration parameters, options and settings for a micro-service
type config struct {
	ServiceName         string                    `yaml:"serviceName"`
	ServiceType         string                    `yaml:"serviceType"`
	HTTPort             int                       `yaml:"httpPort"`
	GRPCPort            int                       `yaml:"grpcPort"`
	HttpOtions          *httpOptions              `yaml:"httpOptions"`
	StartupSleepSeconds int                       `yaml:"startupSleepSeconds"`
	LogLevel            int                       `yaml:"logLevel"`
	Security            *securityOptions          `yaml:"security"`
	Databases           []*databaseOptions        `yaml:"databases"`
	ExternalServices    []*externalServiceOptions `yaml:"externalServices"`
}

// Config contains configuration parameters, options and settings for a micro-service
type Config struct {
	config
}

// New creates config by reading from first non-empty file specified in configFile argument or with --config-file flag
func New(configFile ...string) (*Config, error) {
	cfg := newConfig()

	// Parse the config file
	err := cfg.parse(configFile...)
	if err != nil {
		return nil, err
	}

	// Validate config
	err = errors.Wrap(cfg.validate(), "validation error")
	if err != nil {
		return nil, err
	}

	return &Config{*cfg}, nil
}

const unknownLevel = 1000

func newConfig() *config {
	return &config{
		LogLevel:         unknownLevel,
		Security:         new(securityOptions),
		HttpOtions:       new(httpOptions),
		Databases:        make([]*databaseOptions, 0),
		ExternalServices: make([]*externalServiceOptions, 0),
	}
}
