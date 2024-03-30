package handler

import (
	"database/sql"
)

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) handler {
	return handler{
		db: db,
	}
}
