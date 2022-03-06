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
	fmt.Fprintln(w, "Create a new movie")
}

// TODO: Add a showMovieHandler for the "GET/v1/movies/:id" endpoint.
// For now we retrieve the interpolated "id" parameter from the current URL
// and include it in a placeholder response
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	movie := data.Movie{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Casablanca",
		Runtime: 488,
		Genres: []string{"darma", "romance", "war"},
		Verson: 1,
	}
	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}
