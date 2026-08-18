package core_mysql_conn

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"HOST"     default:"tasks-mysql"`
	Port     string        `envconfig:"PORT"     default:"3306"`
	User     string        `envconfig:"USER"     required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	Database string        `envconfig:"DATABASE" required:"true"`
	Timeout  time.Duration `envconfig:"TIMEOUT"  default:"5s"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("MYSQL", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get MySQL connection pool config: %w", err)
		panic(err)
	}
	return config
}
