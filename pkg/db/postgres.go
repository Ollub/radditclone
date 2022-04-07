package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang-stepik-2022q1/reditclone/config"
	"golang-stepik-2022q1/reditclone/pkg/log"
)

var initial = `
	DROP TABLE IF EXISTS users;
	CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
			pass_hash TEXT NOT NULL
	)
`

func GetPostgres() *sql.DB {
	dsn := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		config.Cfg.DbUser,
		config.Cfg.DbName,
		config.Cfg.DbPass,
		config.Cfg.DbHost,
		config.Cfg.DbPort,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("cant parse config", log.Fields{"error": err})
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Error("ping failed", log.Fields{"error": err})
	}
	db.SetMaxOpenConns(10)

	_, err = db.Exec(initial)
	if err != nil {
		panic(err.Error())
	}
	return db
}
