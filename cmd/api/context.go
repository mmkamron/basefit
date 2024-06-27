package main

import (
	"context"
	"net/http"

	"github.com/mmkamron/basefit/internal/data"
)

type contextKey string

const trainerContextKey = contextKey("trainer")

func (app *application) contextSetTrainer(r *http.Request, trainer *data.Trainer) *http.Request {
	ctx := context.WithValue(r.Context(), trainerContextKey, trainer)
	return r.WithContext(ctx)
}

func (app *application) contextGetTrainer(r *http.Request) *data.Trainer {
	trainer, ok := r.Context().Value(trainerContextKey).(*data.Trainer)
	if !ok {
		panic("missing trainer value in request context")
	}

	return trainer
}
