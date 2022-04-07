package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang-stepik-2022q1/reditclone/config"
	"golang-stepik-2022q1/reditclone/pkg/log"
)

func NewRedis() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Cfg.RedisHost, config.Cfg.RedisPort),
		Password: config.Cfg.RedisPwd, // no password set
		DB:       config.Cfg.RedisDb,  // use default DB
	})
	checkConnection(rdb)
	return &RedisClient{rdb}
}

func checkConnection(client *redis.Client) {
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Debug("Connection to Redis established")
}
