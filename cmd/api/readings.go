package main

import (
	"errors"
	"net/http"

	"github.com/TH-takahirohara/reading_record_api/internal/data"
	"github.com/TH-takahirohara/reading_record_api/internal/validator"
)

func (app *application) createReadingsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BookName       string `json:"book_name"`
		BookAuthor     string `json:"book_author"`
		TotalPageCount int    `json:"total_page_count"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	reading := &data.Reading{
		BookName:       input.BookName,
		BookAuthor:     input.BookAuthor,
		TotalPageCount: input.TotalPageCount,
		UserID:         user.ID,
	}

	v := validator.New()

	if data.ValidateReading(v, reading); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Readings.Insert(reading)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"reading": reading}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showReadingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	reading, err := app.models.Readings.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrNotPermitted):
			app.notPermittedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"reading": reading}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateReadingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	reading, err := app.models.Readings.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrNotPermitted):
			app.notPermittedResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		BookName   *string `json:"book_name"`
		BookAuthor *string `json:"book_author"`
		Memo       *string `json:"memo"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.BookName != nil {
		reading.BookName = *input.BookName
	}
	if input.BookAuthor != nil {
		reading.BookAuthor = *input.BookAuthor
	}
	if input.Memo != nil {
		reading.Memo = *input.Memo
	}

	err = app.models.Readings.Update(reading)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"reading": reading}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteReadingHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Readings.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listReadingsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	readings, err := app.models.Readings.GetAll(user.ID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"readings": readings}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
