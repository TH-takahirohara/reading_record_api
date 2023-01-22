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

	message := "サーバー側の処理において問題が発生しました。"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "更新処理の競合が発生したため、更新に失敗しました。再度リクエストを実行してください"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "無効な認証情報です"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "認証トークンが含まれていないか、無効な値です"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
