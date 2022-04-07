package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type IRedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) IRedisStatusCmd
	Get(ctx context.Context, key string) IRedisStatusCmd
}

type IRedisStatusCmd interface {
	Err() error
	Result() (string, error)
}

type RedisClient struct {
	cli *redis.Client
}

func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) IRedisStatusCmd {
	cmd := rc.cli.Set(ctx, key, value, expiration)
	return &RedisStatusCmd{cmd}
}

func (rc *RedisClient) Get(ctx context.Context, key string) IRedisStatusCmd {
	cmd := rc.cli.Get(ctx, key)
	return &RedisStringCmd{cmd}
}

type RedisStatusCmd struct {
	cmd *redis.StatusCmd
}

func (rs *RedisStatusCmd) Err() error {
	return rs.cmd.Err()
}

func (rs *RedisStatusCmd) Result() (string, error) {
	return rs.cmd.Result()
}

type RedisStringCmd struct {
	cmd *redis.StringCmd
}

func (rs *RedisStringCmd) Err() error {
	return rs.cmd.Err()
}

func (rs *RedisStringCmd) Result() (string, error) {
	return rs.cmd.Result()
}
