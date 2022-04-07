package repo

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"golang-stepik-2022q1/reditclone/pkg/db"
	"golang-stepik-2022q1/reditclone/pkg/session"
	"testing"
	"time"
)

func TestRedisRepo_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	clientMock := db.NewMockIRedisClient(ctrl)
	cmdMock := db.NewMockIRedisStatusCmd(ctrl)

	repo := NewRedis(clientMock)

	sessionId := session.SessionId("sessionID")
	sessionKey := sessionKey(sessionId)
	ttl := time.Duration(0)
	val := ""
	unexpectedErr := errors.New("Unexpected error")

	for _, tt := range [...]struct {
		name        string
		setup       func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd)
		expectedErr error
	}{
		{
			name: "OK",
			setup: func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd) {
				gomock.InOrder(
					clientMock.EXPECT().Set(ctx, sessionKey, val, ttl).Return(cmdMock),
					cmdMock.EXPECT().Err().Return(nil),
				)
			},
			expectedErr: nil,
		},
		{
			name: "Redis error",
			setup: func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd) {
				gomock.InOrder(
					clientMock.EXPECT().Set(ctx, sessionKey, val, ttl).Return(cmdMock),
					cmdMock.EXPECT().Err().Return(unexpectedErr),
				)
			},
			expectedErr: unexpectedErr,
		},
	} {
		tt.setup(clientMock, cmdMock)
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Set(ctx, sessionId)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestRedisRepo_CheckExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	clientMock := db.NewMockIRedisClient(ctrl)
	cmdMock := db.NewMockIRedisStatusCmd(ctrl)

	repo := NewRedis(clientMock)

	sessionId := session.SessionId("sessionID")
	sessionKey := sessionKey(sessionId)
	val := ""
	unexpectedErr := errors.New("Unexpected error")

	for _, tt := range [...]struct {
		name     string
		setup    func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd)
		expected bool
	}{
		{
			name: "OK",
			setup: func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd) {
				gomock.InOrder(
					clientMock.EXPECT().Get(ctx, sessionKey).Return(cmdMock),
					cmdMock.EXPECT().Result().Return(val, nil),
				)
			},
			expected: true,
		},
		{
			name: "Redis error",
			setup: func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd) {
				gomock.InOrder(
					clientMock.EXPECT().Get(ctx, sessionKey).Return(cmdMock),
					cmdMock.EXPECT().Result().Return(val, unexpectedErr),
				)
			},
			expected: false,
		},
		{
			name: "Not found",
			setup: func(cli *db.MockIRedisClient, cmd *db.MockIRedisStatusCmd) {
				gomock.InOrder(
					clientMock.EXPECT().Get(ctx, sessionKey).Return(cmdMock),
					cmdMock.EXPECT().Result().Return(val, redis.ErrNil),
				)
			},
			expected: false,
		},
	} {
		tt.setup(clientMock, cmdMock)
		t.Run(tt.name, func(t *testing.T) {
			res := repo.CheckExists(ctx, sessionId)
			assert.Equal(t, tt.expected, res)
		})
	}
}
