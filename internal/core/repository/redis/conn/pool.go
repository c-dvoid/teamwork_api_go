package core_redis_conn

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Pool struct {
	*redis.Client
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),
	})

	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Pool{Client: client}, nil
}

func NewPoolMust(ctx context.Context, config Config) *Pool {
	pool, err := NewPool(ctx, config)
	if err != nil {
		panic(fmt.Errorf("get Redis connection pool: %w", err))
	}
	return pool
}
