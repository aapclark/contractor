package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*AppConfig, error) {
	var cfg AppConfig

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
