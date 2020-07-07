package config

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// parses config parameters either from environment variables, files or cmd flags
func (cfg *config) parse(froms []configFrom) error {

	from := flag.String(
		"config-from", FromFile.String(),
		`Where to parse config parameters. Options are 'all', 'flags', 'environment' or 'file'`,
	)

	configFile := flag.String(
		"config-file", "config/config.yml",
		`File to read service config`,
	)

	flagCfg := newConfig()

	// calls flag.Parse()
	flagCfg.setConfigFromFlag()

	if len(froms) == 0 {
		froms = []configFrom{fromString(*from)}
	}

	// for remove duplicates
	av := make(map[configFrom]struct{}, 0)

	var err error

loop:
	for _, configFrom := range froms {
		// Check whether its repetition
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

	// set config from secret files
	err = cfg.updateConfigSecrets()
	if err != nil {
		return fmt.Errorf("failed to set config from secrets file: %w", err)
	}

	// call validate and return any errors
	return errors.Wrap(cfg.validate(), "validation error")
}
