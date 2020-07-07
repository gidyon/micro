// Package config contains options for passing setup information for a micro-service
package config

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
	Address      string      `yaml:"address"`
	Host         string      `yaml:"host"`
	Port         int         `yaml:"port"`
	User         string      `yaml:"user"`
	Schema       string      `yaml:"schema"`
	Password     string      `yaml:"password"`
	UserFile     string      `yaml:"userFile"`
	SchemaFile   string      `yaml:"schemaFile"`
	PasswordFile string      `yaml:"passwordFile"`
	Metadata     *dbMetadata `yaml:"metadata"`
}

// databases contains information about databases used by service
type databases struct {
	SQLDatabase   *databaseOptions `yaml:"sqlDatabase"`
	RedisDatabase *databaseOptions `yaml:"redisDatabase"`
}

// externalServiceOptions contains information to connect to a remote service
type externalServiceOptions struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Available   bool   `yaml:"required"`
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
	Databases           *databases                `yaml:"databases"`
	ExternalServices    []*externalServiceOptions `yaml:"externalServices"`
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

	return &Config{*cfg}, nil
}

func newConfig() *config {
	return &config{
		Logging:  new(loggingOptions),
		Security: new(securityOptions),
		Databases: &databases{
			SQLDatabase: &databaseOptions{
				Metadata: new(dbMetadata),
			},
			RedisDatabase: &databaseOptions{
				Metadata: new(dbMetadata),
			},
		},
		ExternalServices: make([]*externalServiceOptions, 0),
	}
}
