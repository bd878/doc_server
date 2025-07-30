package config

import (
	"errors"
	"time"
	"fmt"
	"encoding/json"
)

type (
	WebConfig struct {
		Host string   `json:"host"`
		Port string   `json:"port"`
	}

	Duration struct {
		time.Duration
	}

	PGConfig struct {
		Conn string
	}

	AppConfig struct {
		Environment       string          `json:"environment"`
		LogLevel         string          `json:"log_level"`
		PG               PGConfig        `json:"pg"`
		Web              WebConfig       `json:"web"`
		ShutdownTimeout  Duration        `json:"shutdown_timeout"`
	}
)

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s%s", c.Host, c.Port)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
