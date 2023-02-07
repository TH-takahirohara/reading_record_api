package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/TH-takahirohara/reading_record_api/internal/data"
	"github.com/TH-takahirohara/reading_record_api/internal/validator"
)

func (app *application) createDailyProgressHandler(w http.ResponseWriter, r *http.Request) {
	readingID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	reading, err := app.models.Readings.Get(readingID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if reading.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	var input struct {
		ReadDate time.Time `json:"readDate"`
		ReadPage int       `json:"readPage"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dailyProgress := &data.DailyProgress{
		ReadDate:  input.ReadDate,
		ReadPage:  input.ReadPage,
		ReadingID: readingID,
	}

	latestDailyProgress, err := app.models.DailyProgresses.GetLatest(readingID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if latestDailyProgress != nil {
		v.Check(dailyProgress.ReadDate.After(latestDailyProgress.ReadDate), "read_date", "最後に記録された日付より新しい日付を指定する必要があります")
		v.Check(dailyProgress.ReadPage > latestDailyProgress.ReadPage, "read_page", "最後に記録されたページ番号より大きい値を指定する必要があります")
	}
	v.Check(dailyProgress.ReadPage > 0, "read_page", "ページ番号は0より大きい値を指定する必要があります")
	v.Check(dailyProgress.ReadPage <= reading.TotalPageCount, "read_page", "ページ番号は総ページ数以下の値である必要があります")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DailyProgresses.Insert(dailyProgress)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"dailyProgress": dailyProgress}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteDailyProgressHandler(w http.ResponseWriter, r *http.Request) {
	readingID, err := app.readIntParam(r, "reading_id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dailyProgressID, err := app.readIntParam(r, "daily_progress_id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	reading, err := app.models.Readings.Get(readingID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if reading.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.DailyProgresses.Delete(dailyProgressID)
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
