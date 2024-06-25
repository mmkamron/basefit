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
	mux.HandleFunc("PUT /v1/trainers/activated", app.activateUserHandler)
	return mux
}
