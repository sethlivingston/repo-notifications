package server

import (
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/log"
)

// Uses Moov's config library to load and merge configuration files from
// ./config/config.default.yml and the locations specified by environment
// variables APP_CONFIG and APP_CONFIG_SECRET.
func LoadConfig(logger log.Logger) (*Config, error) {
	configService := config.NewService(logger)

	cfg := &Config{}
	if err := configService.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
