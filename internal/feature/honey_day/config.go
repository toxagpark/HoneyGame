package honey_day

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// INTERVAL — период медового дня, формат time.Duration: 24h, 12h30m, 45s.
	INTERVAL time.Duration `envconfig:"INTERVAL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("HONEY_DAY", cfg); err != nil {
		return &Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if cfg.INTERVAL <= 0 {
		return &Config{}, fmt.Errorf("HONEY_DAY_INTERVAL must be positive, like 24h")
	}

	return cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get honey day config: %w", err)
		panic(err)
	}
	return cfg
}
