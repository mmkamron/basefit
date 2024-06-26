package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/trainers", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/trainers/activated", app.activateUserHandler)
	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return mux
}
