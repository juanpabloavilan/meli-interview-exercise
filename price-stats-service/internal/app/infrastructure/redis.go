package infrastructure

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, redisConnString string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisConnString)
	if err != nil {
		return nil, err
	}

	redis := redis.NewClient(opts)
	res := redis.Ping(ctx)
	if err := res.Err(); err != nil {
		return nil, err
	}

	return redis, nil
}
