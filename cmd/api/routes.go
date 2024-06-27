package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	//TODO:add middleware
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		app.writeJSON(w, 200, "Hello", nil)
	})
	mux.HandleFunc("POST /v1/trainers", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/trainers/activated", app.activateUserHandler)
	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.authenticate(mux))
}
