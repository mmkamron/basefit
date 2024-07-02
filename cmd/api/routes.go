package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck",
		app.requireActivatedUser(func(w http.ResponseWriter, r *http.Request) {
			app.writeJSON(w, 200, "Hello", nil)
		}))

	mux.HandleFunc("POST /v1/trainers", app.createTrainerHandler)
	mux.HandleFunc("GET /v1/trainers/{id}", app.showTrainerHandler)
	mux.HandleFunc("GET /v1/trainers", app.showTrainersHandler)
	mux.HandleFunc("PATCH /v1/trainers/{id}", app.updateTrainerHandler)
	mux.HandleFunc("DELETE /v1/trainers/{id}", app.deleteTrainerHandler)

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)
	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.authenticate(mux))
}
