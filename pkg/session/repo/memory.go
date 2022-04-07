package repo

import (
	"golang-stepik-2022q1/reditclone/pkg/session"
	"sync"
)

type MemRepo struct {
	sync.RWMutex
	// mapping of users id on session
	userSession map[session.SessionId]bool
}

func NewMemRepo() *MemRepo {
	userSession := make(map[session.SessionId]bool)
	return &MemRepo{userSession: userSession}
}

func (r *MemRepo) Set(sessionId session.SessionId) error {
	r.Lock()
	defer r.Unlock()
	r.userSession[sessionId] = true
	return nil
}

func (r *MemRepo) CheckExists(sessionId session.SessionId) bool {
	r.RLock()
	defer r.RUnlock()
	return r.userSession[sessionId]
}
