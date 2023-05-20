package pkg

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB() *sql.DB {
	config := Load()
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
