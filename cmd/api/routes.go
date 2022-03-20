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
	
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMovieHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	return router
}

