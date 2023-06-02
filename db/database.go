package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./basefit.db")
	if err != nil {
		log.Println(err)
	}
	return db
}
