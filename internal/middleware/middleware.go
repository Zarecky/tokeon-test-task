package middleware

import (
	"tokeon-test-task/internal/config"

	"tokeon-test-task/pkg/log"
)

type Middleware struct {
	logger log.Logger
	config *config.Config
}

func New(logger log.Logger, config *config.Config) *Middleware {
	return &Middleware{
		logger: logger,
		config: config,
	}
}
