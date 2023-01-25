package main

import (
	"errors"
	"net/http"

	"github.com/TH-takahirohara/reading_record_api/internal/data"
	"github.com/TH-takahirohara/reading_record_api/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *application) createReadingsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"reading": reading}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
