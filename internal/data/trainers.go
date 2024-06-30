package data

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Trainer struct {
	ID         int64
	Email      string
	Name       string
	Experience int16
	Activities []string
}

type TrainerModel struct {
	DB *sql.DB
}

func (t TrainerModel) Insert(trainer *Trainer) error {
	query := `
        INSERT INTO trainers (email, name, experience, activities)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	args := []interface{}{trainer.Email, trainer.Name, trainer.Experience, pq.Array(trainer.Activities)}

	return t.DB.QueryRow(query, args...).Scan(&trainer.ID)
}

func (t TrainerModel) Get(id int64) (*Trainer, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, email, name, experience, activities
		FROM trainers
		WHERE id = $1`

	var trainer Trainer

	err := t.DB.QueryRow(query, id).Scan(
		&trainer.ID,
		&trainer.Email,
		&trainer.Name,
		&trainer.Experience,
		pq.Array(&trainer.Activities),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &trainer, nil
}

func (t TrainerModel) GetAll() ([]*Trainer, error) {
	rows, err := t.DB.Query("SELECT * FROM trainers")
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	trainers := make([]*Trainer, 0)
	for rows.Next() {
		trainer := new(Trainer)
		err := rows.Scan(
			&trainer.ID,
			&trainer.Email,
			&trainer.Name,
			&trainer.Experience,
			pq.Array(&trainer.Activities),
		)
		if err != nil {
			return nil, err
		}
		trainers = append(trainers, trainer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trainers, nil
}

func (t TrainerModel) Update(trainer *Trainer) error {
	query := `
		UPDATE trainers
		SET email = $1, name = $2, experience = $3, activities = $4
		WHERE id = $5
		RETURNING id
	`

	args := []interface{}{
		trainer.Email,
		trainer.Name,
		trainer.Experience,
		pq.Array(trainer.Activities),
		trainer.ID,
	}

	return t.DB.QueryRow(query, args...).Scan(&trainer.ID)
}

func (t TrainerModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM trainers
		WHERE id = $1
	`

	result, err := t.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
