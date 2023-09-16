package cache

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	GetFromJSON(ctx context.Context, key string, dst any) error
	MGet(ctx context.Context, keys ...string) ([][]byte, error)
	HGet(ctx context.Context, key string, field string) ([]byte, error)
	HMGet(ctx context.Context, key string, fields ...string) ([][]byte, error)
	HGetFromJSON(ctx context.Context, key string, field string, dst any) error
	HGetAll(ctx context.Context, key string) (map[string][]byte, error)
	HGetAllFromJSON(ctx context.Context, key string, dst map[string]any) error
	HMGetFromJSON(ctx context.Context, key string, fields ...string) ([]any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error
	HSet(ctx context.Context, key string, value map[string]any) error
	HSetJSON(ctx context.Context, key string, value map[string]any) error
	BulkSet(ctx context.Context, data map[string]any) error
	BulkSetJSON(ctx context.Context, data map[string]any) error
	Delete(ctx context.Context, key string) error
	HDelete(ctx context.Context, key string, fields ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

type cache struct {
	conn *redis.Client
}

func New(redisClient *redis.Client) Cache {
	return &cache{redisClient}
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := c.conn.Get(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, cmd.Err()
	}
	return cmd.Bytes()
}

func (c *cache) GetFromJSON(ctx context.Context, key string, dst any) error {
	cmd := c.conn.Get(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return ErrNotFound
		}
		return cmd.Err()
	}

	res, err := cmd.Bytes()
	if err != nil {
		return errors.Wrap(err, "failed to get result from cmd")
	}

	if err := json.Unmarshal(res, dst); err != nil {
		return errors.Wrap(err, "failed to unmarshal to dst")
	}

	return nil
}

func (c *cache) MGet(ctx context.Context, keys ...string) ([][]byte, error) {
	cmd := c.conn.MGet(ctx, keys...)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, cmd.Err()
	}

	resAny, err := cmd.Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result from cmd")
	}

	var resBytes [][]byte
	for _, v := range resAny {
		if v == nil {
			resBytes = append(resBytes, nil)
		} else {
			resBytes = append(resBytes, []byte(v.(string)))
		}
	}

	return resBytes, nil
}

func (c *cache) HGet(ctx context.Context, key string, field string) ([]byte, error) {
	cmd := c.conn.HGet(ctx, key, field)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, cmd.Err()
	}
	return cmd.Bytes()
}

func (c *cache) HMGet(ctx context.Context, key string, fields ...string) ([][]byte, error) {
	cmd := c.conn.HMGet(ctx, key, fields...)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, cmd.Err()
	}

	resAny, err := cmd.Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result from cmd")
	}

	resBytes := make([][]byte, len(resAny))
	for k, v := range resAny {
		if v == nil {
			continue
		}
		resBytes[k] = []byte(v.(string))
	}

	return resBytes, nil
}

func (c *cache) HMGetFromJSON(ctx context.Context, key string, fields ...string) ([]any, error) {
	cmd := c.conn.HMGet(ctx, key, fields...)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, cmd.Err()
	}

	resAny, err := cmd.Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result from cmd")
	}

	result := make([]any, len(resAny))
	for k, v := range resAny {
		var res any
		if err := json.Unmarshal([]byte(v.(string)), &res); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal to dst")
		}
		result[k] = res
	}

	return result, nil
}

func (c *cache) HGetFromJSON(ctx context.Context, key string, field string, dst any) error {
	cmd := c.conn.HGet(ctx, key, field)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return ErrNotFound
		}
		return cmd.Err()
	}

	res, err := cmd.Bytes()
	if err != nil {
		return errors.Wrap(err, "failed to get result from cmd")
	}

	if err := json.Unmarshal(res, dst); err != nil {
		return errors.Wrap(err, "failed to unmarshal to dst")
	}

	return nil
}

func (c *cache) HGetAll(ctx context.Context, key string) (map[string][]byte, error) {
	cmd := c.conn.HGetAll(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, cmd.Err()
	}

	resAny, err := cmd.Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get result from cmd")
	}

	resBytes := make(map[string][]byte)
	for k, v := range resAny {
		resBytes[k] = []byte(v)
	}

	return resBytes, nil
}

func (c *cache) HGetAllFromJSON(ctx context.Context, key string, dst map[string]any) error {
	cmd := c.conn.HGetAll(ctx, key)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return ErrNotFound
		}
		return cmd.Err()
	}

	resAny, err := cmd.Result()
	if err != nil {
		return errors.Wrap(err, "failed to get result from cmd")
	}

	for k, v := range resAny {
		var res any
		if err := json.Unmarshal([]byte(v), &res); err != nil {
			return errors.Wrap(err, "failed to unmarshal to dst")
		}
		dst[k] = res
	}

	return nil
}

func (c *cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.conn.Set(ctx, key, value, expiration).Err()
}

func (c *cache) SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	marshalledData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.conn.Set(ctx, key, marshalledData, expiration).Err()
}

func (c *cache) HSet(ctx context.Context, key string, value map[string]any) error {
	return c.conn.HSet(ctx, key, value).Err()
}

func (c *cache) HSetJSON(ctx context.Context, key string, value map[string]any) error {
	marshalledData := make(map[string]any)

	for k, v := range value {
		marshalledValue, err := json.Marshal(v)
		if err != nil {
			return err
		}
		marshalledData[k] = marshalledValue
	}

	return c.conn.HSet(ctx, key, marshalledData).Err()
}

func (c *cache) BulkSet(ctx context.Context, data map[string]any) error {
	return c.conn.MSet(ctx, data).Err()
}

func (c *cache) BulkSetJSON(ctx context.Context, data map[string]any) error {
	marshalledData := make(map[string]string)
	for k, v := range data {
		t, err := json.Marshal(v)
		if err != nil {
			return err
		}
		marshalledData[k] = string(t)
	}

	return c.conn.MSet(ctx, marshalledData).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.conn.Del(ctx, key).Err()
}

func (c *cache) HDelete(ctx context.Context, key string, fields ...string) error {
	return c.conn.HDel(ctx, key, fields...).Err()
}

func (c *cache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.Get(ctx, key)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, ErrNotFound) {
		return false, nil
	}

	return false, err
}

func (c *cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.conn.Expire(ctx, key, expiration).Err()
}
