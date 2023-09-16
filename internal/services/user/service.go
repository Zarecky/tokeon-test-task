package user

import (
	"context"
	"tokeon-test-task/internal/config"
	"tokeon-test-task/internal/models"
	"tokeon-test-task/internal/repos/users"
	"tokeon-test-task/pkg/cache"
	"tokeon-test-task/pkg/log"

	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	Delete(ctx context.Context, user *models.User) (*models.User, error)
}

type service struct {
	config *config.Config
	logger log.Logger
	repo   users.Repo
	cache  cache.Cache
}

func NewService(config *config.Config, logger log.Logger, repo users.Repo, cache cache.Cache) Service {
	return &service{config: config, logger: logger, repo: repo, cache: cache}
}
