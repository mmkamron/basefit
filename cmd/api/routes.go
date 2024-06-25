package main

import (
	"net/http"

	"github.com/mmkamron/basefit/internal/oauth"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	//TODO:group in v1 path
	mux.HandleFunc("GET /v1/auth/google/login", oauth.OauthGoogleLogin)
	mux.HandleFunc("GET /v1/auth/google/callback", oauth.OauthGoogleCallback)

	mux.HandleFunc("POST /v1/trainers", app.registerUserHandler)

	// mux.HandleFunc("GET /trainers", handler.TrainerReadAll)
	// mux.HandleFunc("GET /trainer/{username}", app.showTrainer)
	// mux.HandleFunc("POST /trainer", app.createTrainer)
	// mux.HandleFunc("PUT /trainer/{username}", handler.TrainerUpdate)
	// mux.HandleFunc("DELETE /trainer/{username}", handler.TrainerDelete)

	// mux.HandleFunc("POST /trainer/{username}/schedule", handler.ScheduleCreate)
	// mux.HandleFunc("PUT /trainer/{username}/schedule/{id}", handler.ScheduleUpdate)
	// mux.HandleFunc("DELETE /trainer/{username}/schedule/{id}", handler.ScheduleDelete)
	//
	// mux.HandleFunc("POST /trainer/{username}/booking", handler.BookingCreate)
	// mux.HandleFunc("PUT /trainer/{username}/booking", handler.BookingUpdate)
	// mux.HandleFunc("DELETE /trainer/{username}/booking", handler.BookingUpdate)
	return mux
}
