package config

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// SQLDBType is sql database type
	SQLDBType = "sqlDatabase"
	// RedisDBType is redis database type
	RedisDBType = "redisDatabase"
)

var (
	dbTypes = []string{SQLDBType, RedisDBType}
)

func (cfg *config) validate() error {
	var err error

	// Service section
	switch {
	case strings.TrimSpace(cfg.ServiceName) == "":
		return errors.New("service name is required")
	case cfg.HTTPort == 0:
		return errors.New("service port is required")
	}

	// TLS settings
	if !cfg.Security.Insecure {
		switch {
		case strings.TrimSpace(cfg.Security.TLSCertFile) == "":
			return errors.New("tls cert file is required")
		case strings.TrimSpace(cfg.Security.TLSKeyFile) == "":
			return errors.New("tls key file is required")
		case strings.TrimSpace(cfg.Security.ServerName) == "":
			return errors.New("server name is required")
		}
	}

	// Databases validation
	for _, db := range cfg.Databases {
		err = validateDBOptions(db)
		if err != nil {
			return err
		}
	}

	// External services validation
	for _, srv := range cfg.ExternalServices {
		err = validateService(srv)
		if err != nil {
			return nil
		}
	}

	return nil
}

func validateDBOptions(db *databaseOptions) error {
	switch db.Type {
	case SQLDBType:
	case RedisDBType:
	case "":
		return fmt.Errorf("database type is required. Supported types are %v", dbTypes)
	default:
		return fmt.Errorf("database type %s not known. Supported types are %v", db.Type, dbTypes)
	}

	if db.Required {
		switch {
		case strings.TrimSpace(db.Metadata.Name) == "":
			return errors.New("database name is required")
		case strings.TrimSpace(db.Address) == "":
			return errors.New("database address is required")
		case strings.TrimSpace(db.User) == "" && db.Address == SQLDBType:
			return errors.New("database user is required")
		case strings.TrimSpace(db.Password) == "" && db.Address == SQLDBType:
			return errors.New("database password is required")
		case strings.TrimSpace(db.Schema) == "" && db.Address == SQLDBType:
			return errors.New("database schema is required")
		}
	}

	return nil
}

func validateService(srv *externalServiceOptions) error {
	if !srv.Required {
		return nil
	}

	switch {
	case strings.TrimSpace(srv.Name) == "":
		return errors.New("service name is required")
	case strings.TrimSpace(srv.ServerName) == "" && !srv.Insecure:
		return fmt.Errorf("service %s tls server name is required", strings.ToLower(srv.Name))
	case strings.TrimSpace(srv.TLSCertFile) == "" && !srv.Insecure:
		return fmt.Errorf("service %s tls cert is required", strings.ToLower(srv.Name))
	case strings.TrimSpace(srv.Address) == "":
		return fmt.Errorf("service %s address is required", strings.ToLower(srv.Name))
	}

	return nil
}
