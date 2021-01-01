package config

import (
	"flag"
	"fmt"
)

// parses config parameters either from environment variables, files or cmd flags
func (cfg *config) parse(froms []configFrom) error {
	var err error

	from := flag.String(
		"config-from", FromFile.String(),
		`Where to parse config parameters. Options are 'all', 'flags', 'environment' or 'file'`,
	)

	configFile := flag.String(
		"config-file", "configs/config.yml",
		`File to read service config`,
	)

	flagCfg := newConfig()

	// calls flag.Parse()
	flagCfg.setConfigFromFlag()

	if len(froms) == 0 {
		froms = []configFrom{fromString(*from)}
	}

	// for removing duplicates
	av := make(map[configFrom]struct{}, 0)

loop:
	for _, configFrom := range froms {
		_, ok := av[configFrom]
		if ok {
			continue
		}
		av[configFrom] = struct{}{}

		switch configFrom {
		case FromFlag:
			*cfg = *flagCfg
		case FromEnv:
			err = cfg.setConfigFromEnv()
			if err != nil {
				return fmt.Errorf("failed to set config from flag variables: %w", err)
			}
		case FromFile:
			err = cfg.setConfigFromFile(*configFile)
			if err != nil {
				return err
			}
		default:
			// from flag
			*cfg = *flagCfg

			// from environement
			err = cfg.setConfigFromEnv()
			if err != nil {
				return fmt.Errorf("failed to set config from flag variables: %w", err)
			}

			// from file
			err = cfg.setConfigFromFile(*configFile)
			if err != nil {
				return fmt.Errorf("failed to set config from file: %w", err)
			}

			break loop
		}
	}

	// Merge database to databases slice
	if cfg.Database != nil {
		if cfg.Database.SQLDatabase != nil {
			if cfg.Database.SQLDatabase.Required {
				// Validate db
				err = validateDBOptions(cfg.Database.SQLDatabase)
				if err != nil {
					return err
				}
				// Add to head
				cfg.Databases = append([]*databaseOptions{cfg.Database.SQLDatabase}, cfg.Databases...)
			}
		}
		if cfg.Database.RedisDatabase != nil {
			if cfg.Database.RedisDatabase.Required {
				// Validate db
				err = validateDBOptions(cfg.Database.RedisDatabase)
				if err != nil {
					return err
				}
				// Add to head
				cfg.Databases = append([]*databaseOptions{cfg.Database.RedisDatabase}, cfg.Databases...)
			}
		}
	}

	// update config from secret files
	err = cfg.updateConfigSecrets()
	if err != nil {
		return fmt.Errorf("failed to set config from secrets file: %w", err)
	}

	return nil
}
