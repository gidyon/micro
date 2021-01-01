package config

import (
	"bytes"
	"io/ioutil"

	"github.com/pkg/errors"
)

func (cfg *config) updateConfigWith(newCfg *config) {
	// Set config only if zero value
	cfg.ServiceName = setStringIfEmpty(cfg.ServiceName, newCfg.ServiceName)
	cfg.ServiceType = setStringIfEmpty(cfg.ServiceType, newCfg.ServiceType)
	cfg.HTTPort = setIntIfZero(cfg.HTTPort, newCfg.HTTPort)
	cfg.GRPCPort = setIntIfZero(cfg.GRPCPort, newCfg.GRPCPort)
	cfg.StartupSleepSeconds = setIntIfZero(cfg.StartupSleepSeconds, newCfg.StartupSleepSeconds)

	// Service log
	cfg.LogLevel = setIntIfZero(cfg.LogLevel, newCfg.LogLevel)

	// Service security
	if newCfg.Security != nil {
		cfg.Security.TLSCertFile = setStringIfEmpty(cfg.Security.TLSCertFile, newCfg.Security.TLSCertFile)
		cfg.Security.TLSKeyFile = setStringIfEmpty(cfg.Security.TLSKeyFile, newCfg.Security.TLSKeyFile)
		cfg.Security.ServerName = setStringIfEmpty(cfg.Security.ServerName, newCfg.Security.ServerName)
		cfg.Security.Insecure = setBoolIfEmpty(newCfg.Security.Insecure, newCfg.Security.Insecure)
	}

	isDBNonNil, isRedisNonNil := false, false
	if newCfg.Databases != nil {
		if newCfg.Database.SQLDatabase != nil {
			isDBNonNil = true
		}
		if newCfg.Database.RedisDatabase != nil {
			isRedisNonNil = true
		}
	}

	// SQL Database
	if isDBNonNil {
		cfg.Database.SQLDatabase.Required = setBoolIfEmpty(
			cfg.Database.SQLDatabase.Required, newCfg.Database.SQLDatabase.Required,
		)
		cfg.Database.SQLDatabase.Address = setStringIfEmpty(
			cfg.Database.SQLDatabase.Address, newCfg.Database.SQLDatabase.Address,
		)
		cfg.Database.SQLDatabase.User = setStringIfEmpty(
			cfg.Database.SQLDatabase.User, newCfg.Database.SQLDatabase.User,
		)
		cfg.Database.SQLDatabase.Password = setStringIfEmpty(
			cfg.Database.SQLDatabase.Password, newCfg.Database.SQLDatabase.Password,
		)
		cfg.Database.SQLDatabase.Schema = setStringIfEmpty(
			cfg.Database.SQLDatabase.Schema, newCfg.Database.SQLDatabase.Schema,
		)
		cfg.Database.SQLDatabase.UserFile = setStringIfEmpty(
			cfg.Database.SQLDatabase.UserFile, newCfg.Database.SQLDatabase.UserFile,
		)
		cfg.Database.SQLDatabase.PasswordFile = setStringIfEmpty(
			cfg.Database.SQLDatabase.PasswordFile, newCfg.Database.SQLDatabase.PasswordFile,
		)
		cfg.Database.SQLDatabase.SchemaFile = setStringIfEmpty(
			cfg.Database.SQLDatabase.SchemaFile, newCfg.Database.SQLDatabase.SchemaFile,
		)
		cfg.Database.SQLDatabase.Metadata.Dialect = setStringIfEmpty(
			cfg.Database.SQLDatabase.Metadata.Dialect,
			newCfg.Database.SQLDatabase.Metadata.Dialect,
		)
		cfg.Database.SQLDatabase.Metadata.Name = setStringIfEmpty(
			cfg.Database.SQLDatabase.Metadata.Name,
			newCfg.Database.SQLDatabase.Metadata.Name,
		)
		cfg.Database.SQLDatabase.Metadata.Orm = setStringIfEmpty(
			cfg.Database.SQLDatabase.Metadata.Orm,
			newCfg.Database.SQLDatabase.Metadata.Orm,
		)
	}

	// Redis Database
	if isRedisNonNil {
		cfg.Database.RedisDatabase.Required = setBoolIfEmpty(
			cfg.Database.RedisDatabase.Required, newCfg.Database.RedisDatabase.Required,
		)
		cfg.Database.RedisDatabase.Address = setStringIfEmpty(
			cfg.Database.RedisDatabase.Address, newCfg.Database.RedisDatabase.Address,
		)
		cfg.Database.RedisDatabase.User = setStringIfEmpty(
			cfg.Database.RedisDatabase.User, newCfg.Database.RedisDatabase.User,
		)
		cfg.Database.RedisDatabase.Password = setStringIfEmpty(
			cfg.Database.RedisDatabase.Password, newCfg.Database.RedisDatabase.Password,
		)
		cfg.Database.RedisDatabase.Schema = setStringIfEmpty(
			cfg.Database.RedisDatabase.Schema, newCfg.Database.RedisDatabase.Schema,
		)
		cfg.Database.RedisDatabase.UserFile = setStringIfEmpty(
			cfg.Database.RedisDatabase.UserFile, newCfg.Database.RedisDatabase.UserFile,
		)
		cfg.Database.RedisDatabase.PasswordFile = setStringIfEmpty(
			cfg.Database.RedisDatabase.PasswordFile, newCfg.Database.RedisDatabase.PasswordFile,
		)
		cfg.Database.RedisDatabase.SchemaFile = setStringIfEmpty(
			cfg.Database.RedisDatabase.SchemaFile, newCfg.Database.RedisDatabase.SchemaFile,
		)
		cfg.Database.RedisDatabase.Metadata.Name = setStringIfEmpty(
			cfg.Database.RedisDatabase.Metadata.Name,
			newCfg.Database.RedisDatabase.Metadata.Name,
		)
		cfg.Database.RedisDatabase.Metadata.UseRediSearch = setBoolIfEmpty(
			cfg.Database.RedisDatabase.Metadata.UseRediSearch,
			newCfg.Database.RedisDatabase.Metadata.UseRediSearch,
		)
	}

	// Update databases options
	if len(newCfg.Databases) > 0 {
		cfg.Databases = newCfg.Databases
	}

	// External services
	if len(newCfg.ExternalServices) != 0 {
		// cfg.ExternalServices
		for _, extSrv := range newCfg.ExternalServices {
			cfg.ExternalServices = append(cfg.ExternalServices, extSrv)
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

func setStringFromFileIfEmpty(val, def string) (string, error) {
	if val == "" {
		return getFileContent(def)
	}
	return val, nil
}

func setStringIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func setBoolIfEmpty(val, def bool) bool {
	if val {
		return val
	}
	return def
}

func setSliceIfEmpty(val, def []*externalServiceOptions) []*externalServiceOptions {
	if val == nil {
		if def != nil {
			return def
		}
	}
	return val
}

func setIntIfZero(val, def int) int {
	if val == 0 {
		return def
	}
	return val
}
