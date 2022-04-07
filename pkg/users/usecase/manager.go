package usecase

import (
	"context"
	errors2 "errors"
	"golang-stepik-2022q1/reditclone/pkg/errors"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/users"
	"golang.org/x/crypto/bcrypt"
)

var UserExistsError = errors2.New("UserId exists")

type Repo interface {
	Add(user *users.User) (int64, error)
	GetByName(string) (*users.User, error)
}

type Manager struct {
	repo Repo
}

func NewManager(repo Repo) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) Create(ctx context.Context, in *users.UserIn) (*users.User, error) {
	u, err := m.repo.GetByName(in.Name)
	if err != nil {
		log.Clog(ctx).Error("User repo error", log.Fields{"error": err.Error()})
		return nil, errors.InternalError{err.Error()}
	}
	if u != nil {
		log.Clog(ctx).Info("UserId exist")
		return nil, UserExistsError
	}

	hashPass, _ := HashPass(in.Password)
	u = &users.User{
		Name:     in.Name,
		PassHash: hashPass,
	}
	lastId, err := m.repo.Add(u)
	if err != nil {
		return nil, errors.InternalError{err.Error()}
	}
	u.Id = int(lastId)
	log.Clog(ctx).Info("UserId created", log.Fields{"id": u.Id, "name": u.Name})
	return u, nil
}

func (m *Manager) GetByName(ctx context.Context, name string) (*users.User, error) {
	u, err := m.repo.GetByName(name)
	if err != nil {
		log.Clog(ctx).Error("UserId repo error", log.Fields{"error": err.Error()})
		return nil, errors.InternalError{err.Error()}
	}
	return u, nil
}

func CheckPassword(hashPass, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
	if err != nil {
		return false
	}
	return true
}

func HashPass(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), err
}
