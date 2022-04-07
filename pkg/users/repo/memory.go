package repo

import (
	"golang-stepik-2022q1/reditclone/pkg/users"
	"sync"
)

type MemRepo struct {
	sync.RWMutex
	items []*users.User
}

func NewMemRepo() *MemRepo {
	items := make([]*users.User, 0, 10)
	return &MemRepo{items: items}
}

func (r *MemRepo) Add(u *users.User) (*users.User, error) {
	r.Lock()
	defer r.Unlock()
	r.items = append(r.items, u)
	return u, nil
}

func (r *MemRepo) GetByName(val string) (*users.User, error) {
	r.RLock()
	defer r.RUnlock()

	for _, u := range r.items {
		if u.Name == val {
			return u, nil
		}
	}
	return nil, nil
}
