package config

import (
	"tokeon-test-task/pkg/hc"
	"tokeon-test-task/pkg/postgres"
	"tokeon-test-task/pkg/redisclient"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	EnvCI       string `json:"ENV_CI" default:"local"`
	ApiAddr     string `json:"API_ADDR" default:"localhost:8080"`
	ServiceName string `json:"SERVICE_NAME" default:"tokeon-test-task"`
	Port        int    `json:"PORT" default:"8080"`
	Postgres    postgres.Config
	Redis       redisclient.Config
	HealthCheck hc.Config
}

// Validate config
func (c *Config) Validate() error {
	// Validate postgres
	if err := c.Postgres.Validate(); err != nil {
		return err
	}

	// Validate redis
	if err := c.Redis.Validate(); err != nil {
		return err
	}

	return validation.ValidateStruct(
		c,
		validation.Field(&c.ServiceName, validation.Required),
		validation.Field(&c.Port, validation.Required),
	)
}
