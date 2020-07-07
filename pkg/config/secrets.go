package config

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func (cfg *config) updateConfigSecrets() error {
	var (
		err error
		bs  []byte
	)

	// Database secrets
	// user
	if strings.TrimSpace(cfg.Databases.SQLDatabase.UserFile) != "" {
		bs, err = ioutil.ReadFile(cfg.Databases.SQLDatabase.UserFile)
		if err != nil {
			return errors.Wrap(err, "failed to read database username from file")
		}
		cfg.Databases.SQLDatabase.User = string(bytes.TrimSpace(bs))
	}
	// password
	if strings.TrimSpace(cfg.Databases.SQLDatabase.PasswordFile) != "" {
		bs, err = ioutil.ReadFile(cfg.Databases.SQLDatabase.PasswordFile)
		if err != nil {
			return errors.Wrap(err, "failed to read database password from file")
		}
		cfg.Databases.SQLDatabase.Password = string(bytes.TrimSpace(bs))
	}
	// schema
	if strings.TrimSpace(cfg.Databases.SQLDatabase.SchemaFile) != "" {
		bs, err = ioutil.ReadFile(cfg.Databases.SQLDatabase.SchemaFile)
		if err != nil {
			return errors.Wrap(err, "failed to read database schema from file")
		}
		cfg.Databases.SQLDatabase.Schema = string(bytes.TrimSpace(bs))
	}

	// Redis secrets
	// user
	if strings.TrimSpace(cfg.Databases.RedisDatabase.UserFile) != "" {
		bs, err = ioutil.ReadFile(cfg.Databases.RedisDatabase.UserFile)
		if err != nil {
			return errors.Wrap(err, "failed to read redis username from file")
		}
		cfg.Databases.RedisDatabase.User = string(bytes.TrimSpace(bs))
	}
	// password
	if strings.TrimSpace(cfg.Databases.RedisDatabase.PasswordFile) != "" {
		bs, err = ioutil.ReadFile(cfg.Databases.RedisDatabase.UserFile)
		if err != nil {
			return errors.Wrap(err, "failed to read redis password from file")
		}
		cfg.Databases.RedisDatabase.Password = string(bytes.TrimSpace(bs))
	}

	return nil
}
