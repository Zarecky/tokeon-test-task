package services

import (
	"tokeon-test-task/internal/config"
	"tokeon-test-task/internal/repos"
	"tokeon-test-task/internal/services/user"
	"tokeon-test-task/pkg/cache"
	"tokeon-test-task/pkg/log"
)

type Services interface {
	User() user.Service
}

type services struct {
	userService user.Service
}

func New(logger log.Logger, config *config.Config, repos repos.Repos, cache cache.Cache) (Services, error) {
	return &services{
		userService: user.NewService(config, logger, repos.Users(), cache),
	}, nil
}

func (s *services) User() user.Service {
	return s.userService
}
