package session

import (
	"context"
	"github.com/dgrijalva/jwt-go"
)

type SessionId string

type UserClaims struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
}

type Session struct {
	Id   SessionId  `json:"id"`
	Iat  int64      `json:"iat"`
	Exp  int64      `json:"exp"`
	User UserClaims `json:"users"`
	jwt.StandardClaims
}

func FromCtx(ctx context.Context) *Session {
	session, ok := ctx.Value(SessionKey).(*Session)
	if !ok || session == nil {
		return nil
	}
	return session
}
