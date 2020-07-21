package config

import (
	"errors"
	"fmt"
	"strings"
)

func (cfg *config) validate() error {
	// Service section
	switch {
	case strings.TrimSpace(cfg.ServiceName) == "":
		return errors.New("service name is required")
	case cfg.ServicePort == 0:
		return errors.New("service port is required")
	}

	// Database section
	if cfg.Databases.SQLDatabase.Required {
		switch {
		case strings.TrimSpace(cfg.Databases.SQLDatabase.Host) == "":
			return errors.New("sqldatabase host is required")
		case strings.TrimSpace(cfg.Databases.SQLDatabase.User) == "":
			return errors.New("sqldatabase user is required")
		case strings.TrimSpace(cfg.Databases.SQLDatabase.Password) == "":
			return errors.New("sqldatabase password is required")
		case strings.TrimSpace(cfg.Databases.SQLDatabase.Schema) == "":
			return errors.New("sqldatabase schema is required")
		case cfg.Databases.SQLDatabase.Port == 0:
			cfg.Databases.SQLDatabase.Port = 3306
		}

		if cfg.Databases.SQLDatabase.Address == "" {
			cfg.Databases.SQLDatabase.Address = fmt.Sprintf(
				"%s:%d", cfg.Databases.SQLDatabase.Host, cfg.Databases.SQLDatabase.Port,
			)
		}
	}

	// Redis Section
	if cfg.Databases.RedisDatabase.Required {
		switch {
		case strings.TrimSpace(cfg.Databases.RedisDatabase.Host) == "":
			return errors.New("redis host is required")
		case cfg.Databases.RedisDatabase.Port == 0:
			cfg.Databases.RedisDatabase.Port = 6379
		}

		if cfg.Databases.RedisDatabase.Address == "" {
			cfg.Databases.RedisDatabase.Address = fmt.Sprintf(
				"%s:%d", cfg.Databases.RedisDatabase.Host, cfg.Databases.RedisDatabase.Port,
			)
		}
	}

	// validate tls settings
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

	// validate required external services
	for _, srv := range cfg.ExternalServices {
		if !srv.Available {
			continue
		}
		switch {
		case strings.TrimSpace(srv.Name) == "":
			return errors.New("service name is required")
		case strings.TrimSpace(srv.ServerName) == "" && !srv.Insecure:
			return fmt.Errorf("service %s tls server name is required", strings.ToLower(srv.Name))
		case strings.TrimSpace(srv.TLSCertFile) == "" && !srv.Insecure:
			return fmt.Errorf("service %s tls cert is required", strings.ToLower(srv.Name))
		}
		if strings.TrimSpace(srv.Address) == "" {
			switch {
			case strings.TrimSpace(srv.Host) == "":
				return fmt.Errorf("service %s host is required", strings.ToLower(srv.Name))
			case srv.Port == 0:
				return fmt.Errorf("service %s port is required", strings.ToLower(srv.Name))
			}
		}
	}

	return nil
}
