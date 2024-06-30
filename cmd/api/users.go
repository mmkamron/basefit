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

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	if err := user.Password.Set(input.Password); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "could not process your request", http.StatusBadRequest)
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		http.Error(w, "could not validate the data", http.StatusUnprocessableEntity)
		return
	}

	err := app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			http.Error(w, "a user with this email address already exists", http.StatusUnprocessableEntity)
			return
		default:
			app.logger.Error(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
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
			"ID":              user.ID,
		}

		err := app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	}()

	err = app.writeJSON(w, http.StatusAccepted, user, nil)
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

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
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

	user.Activated = true

	err = app.models.Users.Update(user)
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

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = app.writeJSON(w, http.StatusOK, user, nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
