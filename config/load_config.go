package config

import (
	"os"
	"time"
	"encoding/json"
)

func LoadConfig(configPath string) AppConfig {
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cfg := AppConfig{
		Environment: "development",
		LogLevel:    "DEBUG",
		Web: WebConfig{
			Host: "0.0.0.0",
			Port: "8080",
		},
		ShutdownTimeout: Duration{
			Duration: time.Duration(1*time.Minute),
		},
	}

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}

	return cfg
}