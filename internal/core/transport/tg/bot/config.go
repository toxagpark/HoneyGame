package tgbot

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CHAT_ID int64 `envconfig:"CHAT_ID"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("TG", cfg); err != nil {
		return &Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if cfg.CHAT_ID == 0 {
		return &Config{}, fmt.Errorf("TG_CHAT_ID is not set")
	}

	return cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get tg config: %w", err)
		panic(err)
	}
	return cfg
}
