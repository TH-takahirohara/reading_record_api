package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.GET("/v1/healthcheck", app.healthcheckHandler)

	router.POST("/v1/users", app.registerUserHandler)
	router.PUT("/v1/users/activated", app.activateUserHandler)

	router.POST("/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.logRequest(app.secureHeaders(app.authenticate(router))))
}
