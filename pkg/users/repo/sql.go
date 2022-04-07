package repo

import (
	"database/sql"
	"golang-stepik-2022q1/reditclone/pkg/users"
)

type RepoSql struct {
	db *sql.DB
}

func NewSql(db *sql.DB) *RepoSql {
	return &RepoSql{db: db}
}

func (repo *RepoSql) GetByName(name string) (*users.User, error) {
	user := &users.User{}

	err := repo.db.
		QueryRow(`SELECT id, name, pass_hash FROM users WHERE name = $1`, name).
		Scan(&user.Id, &user.Name, &user.PassHash)
	if err == sql.ErrNoRows {
		// users not found - it's not an error
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *RepoSql) Add(u *users.User) (int64, error) {
	var lastInsertId int64
	err := repo.db.QueryRow(
		`INSERT INTO users ("name", "pass_hash") VALUES ($1, $2) RETURNING id`,
		u.Name,
		u.PassHash,
	).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}
