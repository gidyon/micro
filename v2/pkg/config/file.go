package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (cfg *config) setConfigFromFile(filename string) error {
	if filename == "" {
		filename = "configs/config.yml"
	}
	cfgFromFile, err := readFromYAML(filename)
	if err != nil {
		return errors.Wrap(err, "failed to read config from yaml file")
	}

	// Update config with config from file
	cfg.updateConfigWith(cfgFromFile)

	return nil
}

func readFromYAML(filename string) (*config, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from file")
	}

	cfg := newConfig()

	err = yaml.UnmarshalStrict(bs, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	return cfg, nil
}
