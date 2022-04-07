package usecase

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang-stepik-2022q1/reditclone/config"
	"golang-stepik-2022q1/reditclone/pkg/errors"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/session"
	"golang-stepik-2022q1/reditclone/pkg/users"
	"time"
)

type Repo interface {
	Set(context.Context, session.SessionId) error
	CheckExists(context.Context, session.SessionId) bool
}

type Manager struct {
	repo Repo
}

func NewManager(repo Repo) *Manager {
	return &Manager{repo}
}

func (m *Manager) IssueToken(ctx context.Context, u *users.User) (string, error) {
	sess := &session.Session{
		Id:   session.SessionId(uuid.New().String()),
		User: session.UserClaims{u.Name, u.Id},
		Iat:  time.Now().Unix(),
		Exp:  expDate().Unix(),
	}
	err := m.repo.Set(ctx, sess.Id)
	if err != nil {
		log.Clog(ctx).Error("Error during session creation", log.Fields{"error": err.Error()})
		return "", errors.InternalError{"Error during session creation"}
	}
	tokenString, err := generateToken(sess)
	if err != nil {
		detail := "Error during jwt token generation"
		log.Clog(ctx).Error(detail, log.Fields{"error": err.Error()})
		return "", errors.InternalError{detail}
	}
	return tokenString, nil
}

func (m *Manager) Check(ctx context.Context, token string) (*session.Session, error) {
	sess, err := loadSession(token)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if sess.Exp < now.Unix() {
		return nil, ExpiredTokenErr
	}

	ok := m.repo.CheckExists(ctx, sess.Id)
	if !ok {
		return sess, SessionNotFound
	}

	return sess, nil
}

func expDate() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}

func loadSession(token string) (*session.Session, error) {
	sess := &session.Session{}
	tkn, err := jwt.ParseWithClaims(token, sess, func(token *jwt.Token) (interface{}, error) {
		return config.Cfg.JwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, InvalidTokenErr
	}
	return sess, nil
}

func generateToken(sess *session.Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sess)
	return token.SignedString(config.Cfg.JwtKey)
}
