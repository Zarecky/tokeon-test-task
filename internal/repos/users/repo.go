package users

import "tokeon-test-task/pkg/cache"

type Repo interface {
	Querier
}

type repo struct {
	*Queries
	psql  DBTX
	cache cache.Cache
}

func NewRepo(psql DBTX, cache cache.Cache) Repo {
	return &repo{
		Queries: New(psql),
		psql:    psql,
		cache:   cache,
	}
}
