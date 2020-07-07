package config

import (
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"
)

func (cfg *config) updateConfigWith(newCfg *config) {
	// Set config only if zero value
	cfg.ServiceVersion = setStringIfEmpty(cfg.ServiceVersion, newCfg.ServiceVersion)
	cfg.ServiceName = setStringIfEmpty(cfg.ServiceName, newCfg.ServiceName)
	cfg.ServicePort = setIntIfZero(cfg.ServicePort, newCfg.ServicePort)
	cfg.StartupSleepSeconds = setIntIfZero(cfg.StartupSleepSeconds, newCfg.StartupSleepSeconds)

	// Service logging
	if newCfg.Logging != nil {
		cfg.Logging.Level = setIntIfZero(cfg.Logging.Level, newCfg.Logging.Level)
		cfg.Logging.TimeFormat = setStringIfEmpty(cfg.Logging.TimeFormat, newCfg.Logging.TimeFormat)
	}

	// Service security
	if newCfg.Security != nil {
		cfg.Security.TLSCertFile = setStringIfEmpty(cfg.Security.TLSCertFile, newCfg.Security.TLSCertFile)
		cfg.Security.TLSKeyFile = setStringIfEmpty(cfg.Security.TLSKeyFile, newCfg.Security.TLSKeyFile)
		cfg.Security.ServerName = setStringIfEmpty(cfg.Security.ServerName, newCfg.Security.ServerName)
		cfg.Security.Insecure = setBoolIfEmpty(newCfg.Security.Insecure, newCfg.Security.Insecure)
	}

	isDBNonNil, isRedisNonNil := false, false
	if newCfg.Databases != nil {
		if newCfg.Databases.SQLDatabase != nil {
			isDBNonNil = true
		}
		if newCfg.Databases.RedisDatabase != nil {
			isRedisNonNil = true
		}
	}

	// SQL Database
	if isDBNonNil {
		cfg.Databases.SQLDatabase.Required = setBoolIfEmpty(
			cfg.Databases.SQLDatabase.Required, newCfg.Databases.SQLDatabase.Required,
		)
		cfg.Databases.SQLDatabase.Address = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Address, newCfg.Databases.SQLDatabase.Address,
		)
		cfg.Databases.SQLDatabase.Host = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Host, newCfg.Databases.SQLDatabase.Host,
		)
		cfg.Databases.SQLDatabase.Port = setIntIfZero(
			cfg.Databases.SQLDatabase.Port, newCfg.Databases.SQLDatabase.Port,
		)
		cfg.Databases.SQLDatabase.User = setStringIfEmpty(
			cfg.Databases.SQLDatabase.User, newCfg.Databases.SQLDatabase.User,
		)
		cfg.Databases.SQLDatabase.Password = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Password, newCfg.Databases.SQLDatabase.Password,
		)
		cfg.Databases.SQLDatabase.Schema = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Schema, newCfg.Databases.SQLDatabase.Schema,
		)
		cfg.Databases.SQLDatabase.UserFile = setStringIfEmpty(
			cfg.Databases.SQLDatabase.UserFile, newCfg.Databases.SQLDatabase.UserFile,
		)
		cfg.Databases.SQLDatabase.PasswordFile = setStringIfEmpty(
			cfg.Databases.SQLDatabase.PasswordFile, newCfg.Databases.SQLDatabase.PasswordFile,
		)
		cfg.Databases.SQLDatabase.SchemaFile = setStringIfEmpty(
			cfg.Databases.SQLDatabase.SchemaFile, newCfg.Databases.SQLDatabase.SchemaFile,
		)
		cfg.Databases.SQLDatabase.Metadata.Dialect = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Metadata.Dialect,
			newCfg.Databases.SQLDatabase.Metadata.Dialect,
		)
		cfg.Databases.SQLDatabase.Metadata.Name = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Metadata.Name,
			newCfg.Databases.SQLDatabase.Metadata.Name,
		)
		cfg.Databases.SQLDatabase.Metadata.Orm = setStringIfEmpty(
			cfg.Databases.SQLDatabase.Metadata.Orm,
			newCfg.Databases.SQLDatabase.Metadata.Orm,
		)
	}

	// Redis Database
	if isRedisNonNil {
		cfg.Databases.RedisDatabase.Required = setBoolIfEmpty(
			cfg.Databases.RedisDatabase.Required, newCfg.Databases.RedisDatabase.Required,
		)
		cfg.Databases.RedisDatabase.Address = setStringIfEmpty(
			cfg.Databases.RedisDatabase.Address, newCfg.Databases.RedisDatabase.Address,
		)
		cfg.Databases.RedisDatabase.Host = setStringIfEmpty(
			cfg.Databases.RedisDatabase.Host, newCfg.Databases.RedisDatabase.Host,
		)
		cfg.Databases.RedisDatabase.Port = setIntIfZero(
			cfg.Databases.RedisDatabase.Port, newCfg.Databases.RedisDatabase.Port,
		)
		cfg.Databases.RedisDatabase.User = setStringIfEmpty(
			cfg.Databases.RedisDatabase.User, newCfg.Databases.RedisDatabase.User,
		)
		cfg.Databases.RedisDatabase.Password = setStringIfEmpty(
			cfg.Databases.RedisDatabase.Password, newCfg.Databases.RedisDatabase.Password,
		)
		cfg.Databases.RedisDatabase.Schema = setStringIfEmpty(
			cfg.Databases.RedisDatabase.Schema, newCfg.Databases.RedisDatabase.Schema,
		)
		cfg.Databases.RedisDatabase.UserFile = setStringIfEmpty(
			cfg.Databases.RedisDatabase.UserFile, newCfg.Databases.RedisDatabase.UserFile,
		)
		cfg.Databases.RedisDatabase.PasswordFile = setStringIfEmpty(
			cfg.Databases.RedisDatabase.PasswordFile, newCfg.Databases.RedisDatabase.PasswordFile,
		)
		cfg.Databases.RedisDatabase.SchemaFile = setStringIfEmpty(
			cfg.Databases.RedisDatabase.SchemaFile, newCfg.Databases.RedisDatabase.SchemaFile,
		)
		cfg.Databases.RedisDatabase.Metadata.Name = setStringIfEmpty(
			cfg.Databases.RedisDatabase.Metadata.Name,
			newCfg.Databases.RedisDatabase.Metadata.Name,
		)
		cfg.Databases.RedisDatabase.Metadata.UseRediSearch = setBoolIfEmpty(
			cfg.Databases.RedisDatabase.Metadata.UseRediSearch,
			newCfg.Databases.RedisDatabase.Metadata.UseRediSearch,
		)
	}

	// External services
	if len(cfg.ExternalServices) == 0 {
		cfg.ExternalServices = make([]*externalServiceOptions, 0)
		if len(newCfg.ExternalServices) != 0 {
			// cfg.ExternalServices
			for _, extSrv := range newCfg.ExternalServices {
				cfg.ExternalServices = append(cfg.ExternalServices, extSrv)
			}
		}
	}
}

func getFileContent(fileFile string) (string, error) {
	bs, err := ioutil.ReadFile(fileFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to read from file")
	}
	return string(bytes.TrimSpace(bs)), nil
}

func setStringFromFileIfEmpty(fromString, fileName string) (string, error) {
	if fromString == "" {
		return getFileContent(fileName)
	}
	return fromString, nil
}

func setStringIfEmpty(fromString, toString string) string {
	if fromString == "" {
		return toString
	}
	return fromString
}

func setBoolIfEmpty(from, to bool) bool {
	if from {
		return from
	}
	return to
}

func setSliceIfEmpty(fromSlice, toSlice []*externalServiceOptions) []*externalServiceOptions {
	if fromSlice == nil {
		if toSlice != nil {
			return toSlice
		}
	}
	return fromSlice
}

func setIntIfZero(from, to int) int {
	if from == 0 {
		return to
	}
	return from
}
