package config

import (
	"flag"
	"fmt"
)

// parses config parameters either from environment variables, files or cmd flags
func (cfg *config) parse(cf ...string) error {
	var err error

	configFile := flag.String(
		"config-file", "configs/config.yml",
		`File location to read config parameter`,
	)

	flag.Parse()

	// Update config
	err = cfg.setConfigFromFile(firstVal(append(cf, *configFile)...))
	if err != nil {
		return err
	}

	// Update config from secret files
	err = cfg.updateConfigSecrets()
	if err != nil {
		return fmt.Errorf("failed to set config from secrets file: %w", err)
	}

	return nil
}

func firstVal(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}
