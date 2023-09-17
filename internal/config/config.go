package config

import (
	"tokeon-test-task/pkg/hc"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	EnvCI       string `json:"ENV_CI" default:"local"`
	ApiAddr     string `json:"API_ADDR" default:"localhost:8080"`
	ServiceName string `json:"SERVICE_NAME" default:"tokeon-test-task"`
	Port        int    `json:"PORT" default:"8080"`
	HealthCheck hc.Config
}

// Validate config
func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.ServiceName, validation.Required),
		validation.Field(&c.Port, validation.Required),
	)
}
