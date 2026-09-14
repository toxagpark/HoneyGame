package core_postgres

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	POSTGRES_USER     string        `envconfig:"USER"`
	POSTGRES_PASSWORD string        `envconfig:"PASSWORD"`
	POSTGRES_HOST     string        `envconfig:"HOST"`
	POSTGRES_DB       string        `envconfig:"DB"`
	POSTGRES_PORT     string        `envconfig:"PORT"`
	POSTGRES_TIMEOUT  time.Duration `envconfig:"TIMEOUT"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("POSTGRES", cfg); err != nil {
		return &Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}
	return cfg
}
