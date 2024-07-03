package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&trainer.ID)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
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

func (t TrainerModel) GetAll(name string, activities []string, filters Filters) ([]*Trainer, error) {
	query := fmt.Sprintf(`
		SELECT id, email, name, experience, activities
		FROM trainers
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (activities @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, name, pq.Array(activities))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trainers := []*Trainer{}

	for rows.Next() {
		var trainer Trainer

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
		trainers = append(trainers, &trainer)
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&trainer.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (t TrainerModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM trainers
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
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
