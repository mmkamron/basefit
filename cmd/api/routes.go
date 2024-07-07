package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		app.writeJSON(w, 200, "Hello", nil)
	})

	mux.HandleFunc("GET /v1/trainers", app.requirePermission("movies:read", app.listTrainersHandler))
	mux.HandleFunc("GET /v1/trainers/{id}", app.requirePermission("movies:read", app.showTrainerHandler))
	mux.HandleFunc("POST /v1/trainers", app.requirePermission("movies:write", app.createTrainerHandler))
	mux.HandleFunc("PATCH /v1/trainers/{id}", app.requirePermission("movies:write", app.updateTrainerHandler))
	mux.HandleFunc("DELETE /v1/trainers/{id}", app.requirePermission("movies:write", app.deleteTrainerHandler))

	mux.HandleFunc("POST /v1/users", app.registerUserHandler)
	mux.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)
	mux.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	//WARNING:secure this endpoint
	mux.Handle("GET /v1/metrics", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCors(app.rateLimit(app.authenticate(mux)))))
}
