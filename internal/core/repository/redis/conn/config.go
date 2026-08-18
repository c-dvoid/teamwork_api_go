package core_redis_conn

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host    string        `envconfig:"HOST"    default:"tasks-redis"`
	Port    string        `envconfig:"PORT"    default:"6379"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"5s"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Redis connection pool config: %w", err)
		panic(err)
	}
	return config
}
