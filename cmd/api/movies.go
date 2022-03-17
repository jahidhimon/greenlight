package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/jahidhimon/greenlight.git/internal/data"
	"github.com/jahidhimon/greenlight.git/internal/validator"
)

func (a *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information of the http request body
	// This struct will be our *target decode destination.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from input to a value of movie struct
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	// Initialize new Validator
	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer
	// to the validated movie struct.
	err = a.models.Movies.Insert(movie)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// Dump the contents of the input struct in a HTTP response
	err = a.writeJson(w, 200, envelop{"movie": movie}, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
	fmt.Fprintf(os.Stdout, "%+v\n", movie)
}

// REVIEW: Add a showMovieHandler for the "GET/v1/movies/:id" endpoint.
// For now we retrieve the interpolated "id" parameter from the current URL
// and include it in a placeholder response
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
