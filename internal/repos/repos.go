package repos

import (
	"tokeon-test-task/internal/repos/users"
	"tokeon-test-task/pkg/cache"
	"tokeon-test-task/pkg/log"
	"tokeon-test-task/pkg/postgres"
)

type Repos interface {
	Users() users.Repo
}

type repos struct {
	users users.Repo
}

func New(logger log.Logger, psql *postgres.PostgreSQL, cache cache.Cache) Repos {
	return &repos{
		users: users.NewRepo(psql, cache),
	}
}

func (r *repos) Users() users.Repo {
	return r.users
}
