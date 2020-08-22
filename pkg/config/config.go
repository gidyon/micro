// Package config contains options for bootstrapping a service dependencies
package config

import "github.com/pkg/errors"

type securityOptions struct {
	TLSCertFile string `yaml:"tlsCert"`
	TLSKeyFile  string `yaml:"tlsKey"`
	ServerName  string `yaml:"serverName"`
	Insecure    bool   `yaml:"insecure"`
}

type loggingOptions struct {
	Disabled   bool   `yaml:"disabled"`
	Level      int    `yaml:"level"`
	TimeFormat string `yaml:"timeFormat"`
}

type dbMetadata struct {
	Name          string `yaml:"name"`
	Dialect       string `yaml:"dialect"`
	Orm           string `yaml:"orm"`
	UseRediSearch bool   `yaml:"useRediSearch"`
}

// databaseOptions contains parameters that open connection to a database
type databaseOptions struct {
	Required     bool        `yaml:"required"`
	Type         string      `yaml:"type"`
	Address      string      `yaml:"address"`
	User         string      `yaml:"user"`
	Schema       string      `yaml:"schema"`
	Password     string      `yaml:"password"`
	UserFile     string      `yaml:"userFile"`
	SchemaFile   string      `yaml:"schemaFile"`
	PasswordFile string      `yaml:"passwordFile"`
	Metadata     *dbMetadata `yaml:"metadata"`
}

// databases contains information about databases used by service
type database struct {
	SQLDatabase   *databaseOptions `yaml:"sqlDatabase"`
	RedisDatabase *databaseOptions `yaml:"redisDatabase"`
}

// externalServiceOptions contains information to connect to a remote service
type externalServiceOptions struct {
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
	K8Service   bool   `yaml:"k8service"`
	Address     string `yaml:"address"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	TLSCertFile string `yaml:"tlsCert"`
	ServerName  string `yaml:"serverName"`
	Insecure    bool   `yaml:"insecure"`
}

// config contains configuration parameters, options and settings for a micro-service
type config struct {
	ServiceVersion      string                    `yaml:"serviceVersion"`
	ServiceName         string                    `yaml:"serviceName"`
	ServicePort         int                       `yaml:"servicePort"`
	StartupSleepSeconds int                       `yaml:"startupSleepSeconds"`
	Logging             *loggingOptions           `yaml:"logging"`
	Security            *securityOptions          `yaml:"security"`
	Database            *database                 `yaml:"database"`
	Databases           []*databaseOptions        `yaml:"databases"`
	ExternalServices    []*externalServiceOptions `yaml:"externalServices"`
	Hello               string                    `yaml:"hello"`
}

// Config contains configuration parameters, options and settings for a micro-service
type Config struct {
	config
}

type configFrom int

func (from configFrom) String() string {
	switch from {
	case FromFile:
		return "FILE"
	case FromEnv:
		return "ENV"
	case FromFlag:
		return "FLAG"
	default:
		return "ALL"
	}
}

func fromString(from string) configFrom {
	switch from {
	case FromFile.String():
		return FromFile
	case FromEnv.String():
		return FromEnv
	case FromFlag.String():
		return FromFlag
	default:
		return FromAll
	}
}

const (
	// FromFile is option to read config from yaml file
	FromFile configFrom = 1
	// FromEnv is option to read config from environment variables
	FromEnv configFrom = 2
	// FromFlag is option to read config from flags
	FromFlag configFrom = 3
	// FromAll is option to read config from flags, environment and file
	FromAll configFrom = 4
)

var allFroms = []configFrom{FromAll, FromFlag, FromEnv, FromFile}

func allowedFrom(from configFrom) bool {
	for _, v := range allFroms {
		if from == v {
			return true
		}
	}
	return false
}

// New creates and parses a new config object
func New(from ...configFrom) (*Config, error) {
	cfg := newConfig()

	// Parse the config
	err := cfg.parse(from)
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

func newConfig() *config {
	return &config{
		Logging:  new(loggingOptions),
		Security: new(securityOptions),
		Database: &database{
			SQLDatabase: &databaseOptions{
				Type: SQLDBType,
				Metadata: &dbMetadata{
					Name: "mysql",
				},
			},
			RedisDatabase: &databaseOptions{
				Type: RedisDBType,
				Metadata: &dbMetadata{
					Name: "redis",
				},
			},
		},
		Databases:        make([]*databaseOptions, 0),
		ExternalServices: make([]*externalServiceOptions, 0),
	}
}
