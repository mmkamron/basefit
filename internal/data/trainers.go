package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/mmkamron/basefit/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type TrainerModel struct {
	DB *sql.DB
}

var AnonymousTrainer = &Trainer{}

type Trainer struct {
	ID        int64
	Email     string
	Password  password `json:"-"`
	Activated bool
	Name      string
}

type password struct {
	plaintext *string
	hash      []byte
}

func (t *Trainer) IsAnonymous() bool {
	return t == AnonymousTrainer
}

func (m TrainerModel) Insert(trainer *Trainer) error {
	query := `
        INSERT INTO trainer (name, email, password_hash, activated) 
        VALUES ($1, $2, $3, $4)
        RETURNING id`
	args := []interface{}{trainer.Name, trainer.Email, trainer.Password.hash, trainer.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&trainer.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "trainer_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m TrainerModel) GetByEmail(email string) (*Trainer, error) {
	query := `
        SELECT id, name, email, password_hash, activated
        FROM trainer
        WHERE email = $1`

	var trainer Trainer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&trainer.ID,
		&trainer.Name,
		&trainer.Email,
		&trainer.Password.hash,
		&trainer.Activated,
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

func (m TrainerModel) Update(trainer *Trainer) error {
	query := `
        UPDATE trainer 
        SET name = $1, email = $2, password_hash = $3, activated = $4
        WHERE id = $5
        RETURNING id`
	args := []interface{}{
		trainer.Name,
		trainer.Email,
		trainer.Password.hash,
		trainer.Activated,
		trainer.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&trainer.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "trainer_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m TrainerModel) GetForToken(tokenScope, tokenPlaintext string) (*Trainer, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
		SELECT trainer.id, trainer.name, trainer.email, trainer.password_hash, trainer.activated
        FROM trainer
        INNER JOIN token
        ON trainer.id = token.trainer_id
        WHERE token.hash = $1
        AND token.scope = $2
        AND token.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var trainer Trainer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&trainer.ID,
		&trainer.Name,
		&trainer.Email,
		&trainer.Password.hash,
		&trainer.Activated,
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

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateTrainer(v *validator.Validator, trainer *Trainer) {
	v.Check(trainer.Name != "", "name", "must be provided")
	v.Check(len(trainer.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, trainer.Email)

	if trainer.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *trainer.Password.plaintext)
	}

	if trainer.Password.hash == nil {
		panic("missing password hash for user")
	}
}
