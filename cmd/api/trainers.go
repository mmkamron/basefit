package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mmkamron/basefit/internal/data"
	"github.com/mmkamron/basefit/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not process your request", http.StatusBadRequest)
		return
	}

	trainer := &data.Trainer{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	if err := trainer.Password.Set(input.Password); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not process your request", http.StatusBadRequest)
	}

	v := validator.New()

	if data.ValidateTrainer(v, trainer); !v.Valid() {
		http.Error(w, "could not validate the data", http.StatusUnprocessableEntity)
		return
	}

	err := app.models.Trainers.Insert(trainer)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a trainer with this email address already exists")
			http.Error(w, "a trainer with this email address already exists", http.StatusUnprocessableEntity)
			return
		default:
			app.logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := app.models.Tokens.New(trainer.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Errorf("%s", err).Error())
			}
		}()

		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"ID":              trainer.ID,
		}

		err := app.mailer.Send(trainer.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	}()

	err = app.writeJSON(w, http.StatusAccepted, trainer, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not process your request", http.StatusBadRequest)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		http.Error(w, "could not validate the data", http.StatusUnprocessableEntity)
		return
	}

	trainer, err := app.models.Trainers.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			http.Error(w, "could not validate the data", http.StatusUnprocessableEntity)
		default:
			app.logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	trainer.Activated = true

	err = app.models.Trainers.Update(trainer)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			http.Error(w, "data conflict", http.StatusConflict)
		default:
			app.logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForTrainer(data.ScopeActivation, trainer.ID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, trainer, nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
