package server

import (
	"context"
	"tokeon-test-task/pkg/postgres"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (s *Server) initDB(ctx context.Context) error {
	// Connect to database
	pg, err := postgres.NewPostgreSQL(ctx, s.logger, s.config.Postgres)
	if err != nil {
		return errors.Wrap(err, "error initializing db")
	}
	s.db = pg

	// Connect to redis
	var redisOptions *redis.Options
	if s.config.Redis.URL != "" {
		redisOptions, err = redis.ParseURL(s.config.Redis.URL)
		if err != nil {
			return errors.Wrap(err, "error parsing redis url")
		}
	} else {
		redisOptions = &redis.Options{
			Addr:     s.config.Redis.Addr,
			Username: s.config.Redis.User,
			Password: s.config.Redis.Pass,
			DB:       s.config.Redis.DbIndex,
		}
	}

	redisConn := redis.NewClient(redisOptions)
	if err := redisConn.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "failed to connect to redis")
	}
	s.redis = redisConn

	s.logger.Info("connected to redis")

	return nil
}
