package bad_bees

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// INTERVAL — период налёта злых пчёл, формат time.Duration: 24h, 12h30m, 45s.
	INTERVAL time.Duration `envconfig:"INTERVAL"`
	// RECENT_PERIOD — «за последнее время»: медведь без боёв за этот период
	// считается АФК и становится целью пчёл.
	RECENT_PERIOD time.Duration `envconfig:"RECENT_PERIOD"`
	// MAX_STING_PERCENT — верхняя граница процента мёда за один укус;
	// конкретный процент разыгрывается от 1 до MAX_STING_PERCENT.
	MAX_STING_PERCENT int `envconfig:"MAX_STING_PERCENT"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := envconfig.Process("BAD_BEES", cfg); err != nil {
		return &Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if cfg.INTERVAL <= 0 {
		return &Config{}, fmt.Errorf("BAD_BEES_INTERVAL must be positive, like 24h")
	}
	if cfg.RECENT_PERIOD <= 0 {
		return &Config{}, fmt.Errorf("BAD_BEES_RECENT_PERIOD must be positive, like 24h")
	}
	if cfg.MAX_STING_PERCENT < 1 || cfg.MAX_STING_PERCENT > 100 {
		return &Config{}, fmt.Errorf("BAD_BEES_MAX_STING_PERCENT must be between 1 and 100")
	}

	return cfg, nil
}

func NewConfigMust() *Config {
	cfg, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get bad bees config: %w", err)
		panic(err)
	}
	return cfg
}
