package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodGet, "/v1/readings", app.requireActivatedUser(app.listReadingsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/readings", app.requireActivatedUser(app.createReadingsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/readings/:id", app.requireActivatedUser(app.showReadingHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/readings/:id", app.requireActivatedUser(app.updateReadingHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/readings/:id", app.requireActivatedUser(app.deleteReadingHandler))

	router.HandlerFunc(http.MethodPost, "/v1/readings/:id/daily_progresses", app.requireActivatedUser(app.createDailyProgressHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/reading/:reading_id/daily_progresses/:daily_progress_id", app.requireActivatedUser(app.deleteDailyProgressHandler))

	return app.recoverPanic(app.enableCORS(app.logRequest(app.secureHeaders(app.authenticate(router)))))
}
