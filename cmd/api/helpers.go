package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelop map[string]interface{}

// Retrieve the "id" URL parameter from the current request context,
// then convert it to an integer and return it. On error return 0, error

func (app *application) readIDParam(r *http.Request) (int64, error) {
	// When httprouter is parsing a request, any interpolated URL parameters  will
	// be stared in the request context. We can use the ParamsFromContext() to
	// retrieve a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" parameter
	// from the slice. In our project all movies will have a unique ID
	// but the value returned by ByName() is always a string. So we try to convert
	// it to a base 10 integer (with a bit size of 64). If the parameter couldn't be
	// converted, of is less than 1, we know the ID is invalid so we return error
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJson(w http.ResponseWriter, status int,
	data envelop, headers http.Header) error {
	// Encode the data to json
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Append a newline to make it nicer in terminals
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	// At this point everything happened correctly
	return nil
}
