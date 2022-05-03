package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"kees/server/config"
	"kees/server/helpers"
)

var DB *sql.DB

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func Configure(config config.DatabaseConfig) error {
	helpers.Debug(config)
	var err error
	DB, err = sql.Open("sqlite3", config.Path)
	if err != nil {
		return err
	}
	helpers.Debug(DB)

	return DB.Ping()
}
