package repo

import (
	"context"
	"golang-stepik-2022q1/reditclone/pkg/db"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/session"
)

type RedisRepo struct {
	client db.IRedisClient
}

func NewRedis(cli db.IRedisClient) *RedisRepo {
	return &RedisRepo{cli}
}

func (r *RedisRepo) Set(ctx context.Context, sessionId session.SessionId) error {
	key := sessionKey(sessionId)
	err := r.client.Set(ctx, key, "", 0).Err()
	if err != nil {
		return err
	}
	log.Clog(ctx).Debug("Session set", log.Fields{"key": key})
	return nil
}

func (r *RedisRepo) CheckExists(ctx context.Context, sessionId session.SessionId) bool {
	key := sessionKey(sessionId)
	_, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Clog(ctx).Debug("Session not found", log.Fields{"key": key, "err": err.Error()})
		return false
	}
	log.Clog(ctx).Debug("Session found", log.Fields{"key": key})
	return true
}

func sessionKey(sessionId session.SessionId) string {
	return "session:" + string(sessionId)
}
