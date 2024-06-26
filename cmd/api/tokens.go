package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mmkamron/basefit/internal/data"
	"github.com/mmkamron/basefit/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not process your request", http.StatusBadRequest)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		http.Error(w, "could not validate the data", http.StatusUnprocessableEntity)
		return
	}

	trainer, err := app.models.Trainers.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	match, err := trainer.Password.Matches(input.Password)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(trainer.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, token, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
