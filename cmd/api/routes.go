package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance
	router := httprouter.New()

	// Convert the notFoundResponse helper method to http handlerFunc and
	// set it as the custom error handler for 404 not found responses
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Convert methodNotAllowedResponse helper to http handlerFunc and set
	// it as the custom error handler for 405 method not allowed
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	
	// Register the relevant methods, URL patterns and handler functions for
	// endpoints using HandlerFunc() method.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc("POST", "/v1/movies", app.createMovieHandler)
	router.HandlerFunc("GET", "/v1/movies/:id", app.showMovieHandler)

	return router
}

