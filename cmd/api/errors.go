package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.WithFields(logrus.Fields{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	}).Error(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "サーバー側での処理において問題が発生しました。"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
