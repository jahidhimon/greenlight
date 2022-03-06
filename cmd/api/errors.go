package main

import (
	"fmt"
	"net/http"
)

// Generic logger for this application.
// TODO: Upgrade this to log request information including http method and URL
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// The errorResponse() method is a generic helper for sending JSON-formatteed error
// messages to the client with a given status code. Note that we're using an interface()
// type for the message parameter, rather than just a string type, as this gives us
// Moer flexibility over the values that we can include in the response
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request,
	status int, message interface{}) {
	env := envelop{"error": message}
	// Write and send the json to the client using writeJSON() helper and then log it
	err := app.writeJson(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// The serverErrorResponse() method will be used when our application encounters an
// runtime Error. It logs the detailed error message, then uses the  errorResponse()
// helper to send a 500 Internal Server Error Code and JSON response to the client
func(app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// The notFoundResponse() method will be used to send a 404 not found status code
// and JSON response to the client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The requested resource could not be found\n")
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// The methodNotAllowedResponse() method will be used to send a 404 not found status code
// and JSON response to the client
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resourse\n", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
} 