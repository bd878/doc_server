package config

import (
	"time"
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type (
	WebConfig struct {
		Host string   `default:"0.0.0.0"`
		Port string   `default:"8080"`
	}

	PGConfig struct {
		Conn string
	}

	AppConfig struct {
		Environment      string
		LogLevel         string          `envconfig:"LOG_LEVEL" default:"DEBUG"`
		PG               PGConfig
		Web              WebConfig
		AdminToken       string          `envconfig:"ADMIN_TOKEN"`
		ShutdownTimeout  time.Duration   `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	}
)

func InitConfig() (cfg AppConfig, err error) {
	err = envconfig.Process("", &cfg)

	return
}

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}