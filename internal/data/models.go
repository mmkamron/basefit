package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users    UserModel
	Trainers TrainerModel
	Tokens   TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Trainers: TrainerModel{DB: db},
		Tokens:   TokenModel{DB: db},
	}
}
