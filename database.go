package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	user     = "postgres"
	password = "library"
	dbname   = "library"
)

func ConnectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
