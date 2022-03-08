package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jahidhimon/greenlight.git/internal/data"
)

// TODO: Add a createMovieHandler for the "POST /v1/movies" endpoint.
// For now we simpley return a plain-text placeholder response
func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information of the http request body
	// This struct will be our *target decode destination.
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use Decode() method to decode the body contents into the input struct.
	// Importantly, notice when we call Decode() we pass a *pointer* to the input
	// Struct as the target decode destination. If there is an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Dump the contents of the input struct in a HTTP response
	// TODO: store the json in file or database
	a.writeJson(w, 200, envelop{"created movie": input}, nil)
	fmt.Println(input)
}

// TODO: Add a showMovieHandler for the "GET/v1/movies/:id" endpoint.
// For now we retrieve the interpolated "id" parameter from the current URL
// and include it in a placeholder response
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   488,
		Genres:    []string{"darma", "romance", "war"},
		Verson:    1,
	}
	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
