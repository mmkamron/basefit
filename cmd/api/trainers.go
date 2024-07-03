package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mmkamron/basefit/internal/data"
	"github.com/mmkamron/basefit/internal/validator"
)

func (app *application) createTrainerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email      string
		Name       string
		Experience int16
		Activities []string
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	trainer := &data.Trainer{
		Email:      input.Email,
		Name:       input.Name,
		Experience: input.Experience,
		Activities: input.Activities,
	}

	//TODO:input validation

	err = app.models.Trainers.Insert(trainer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/trainers/%d", trainer.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"trainer": trainer}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTrainerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	trainer, err := app.models.Trainers.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"trainer": trainer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTrainersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name       string
		Activities []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Activities = app.readCSV(qs, "activities", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 5, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "name", "experience", "-id", "-name", "-experience"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	trainers, err := app.models.Trainers.GetAll(input.Name, input.Activities, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"trainers": trainers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// FIX:fix data race (versioning)
func (app *application) updateTrainerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	trainer, err := app.models.Trainers.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Email      *string
		Name       *string
		Experience *int16
		Activities []string
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Email != nil {
		trainer.Email = *input.Email
	}
	if input.Name != nil {
		trainer.Name = *input.Name
	}
	if input.Experience != nil {
		trainer.Experience = *input.Experience
	}
	if input.Activities != nil {
		trainer.Activities = input.Activities
	}

	//TODO:input validation

	err = app.models.Trainers.Update(trainer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"trainer": trainer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTrainerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Trainers.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "trainer successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
