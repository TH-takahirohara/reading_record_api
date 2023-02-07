package main

import (
	"fmt"
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
	app.errorResponse(w, r, http.StatusInternalServerError, envelope{"msg": message})
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, envelope{"msg": err.Error()})
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "このリソースにアクセスするためには認証が必要です"
	app.errorResponse(w, r, http.StatusUnauthorized, envelope{"msg": message})
}

func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "このリソースにアクセスするためにはユーザーアカウントの有効化が必要です"
	app.errorResponse(w, r, http.StatusForbidden, envelope{"msg": message})
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "このリソースにアクセスするために必要な権限がありません"
	app.errorResponse(w, r, http.StatusForbidden, envelope{"msg": message})
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "要求されたリソースが見つかりませんでした"
	app.errorResponse(w, r, http.StatusNotFound, envelope{"msg": message})
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s メソッドはこのリソースでサポートされていません", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, envelope{"msg": message})
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "更新処理の競合が発生したため、更新に失敗しました。再度リクエストを実行してください"
	app.errorResponse(w, r, http.StatusConflict, envelope{"msg": message})
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "無効な認証情報です"
	app.errorResponse(w, r, http.StatusUnauthorized, envelope{"msg": message})
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "認証トークンが含まれていないか、無効な値です"
	app.errorResponse(w, r, http.StatusUnauthorized, envelope{"msg": message})
}
