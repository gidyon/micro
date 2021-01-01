package config

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func (cfg *config) updateConfigSecrets() error {
	var err error

	// Update db secrets
	for _, db := range cfg.Databases {
		err = updateDatabaseSecret(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateDatabaseSecret(db *databaseOptions) error {
	// user
	if strings.TrimSpace(db.UserFile) != "" {
		bs, err := ioutil.ReadFile(db.UserFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read %s database username from file", db.Type)
		}
		db.User = string(bytes.TrimSpace(bs))
	}

	// password
	if strings.TrimSpace(db.PasswordFile) != "" {
		bs, err := ioutil.ReadFile(db.PasswordFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read %s database password from file", db.Type)
		}
		db.Password = string(bytes.TrimSpace(bs))
	}

	// schema
	if strings.TrimSpace(db.SchemaFile) != "" {
		bs, err := ioutil.ReadFile(db.SchemaFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read %s database schema from file", db.Type)
		}
		db.Schema = string(bytes.TrimSpace(bs))
	}

	return nil
}
