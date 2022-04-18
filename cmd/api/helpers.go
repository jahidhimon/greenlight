package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jahidhimon/greenlight.git/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type envelope map[string]interface{}

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

func (app *application) writeJSON(w http.ResponseWriter, status int,
	data envelope, headers http.Header) error {
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

func (app *application) readJSON(w http.ResponseWriter,
	r *http.Request, dst interface{}) error {
	// Use http.MaxBytesReader to limit the size of the request body to 1MB
	maxbytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxbytes))

	// Initialize the json.Decoder and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes
	// any field which cannot be mapped to the target destination, the decoder will
	// return an error instead of just ignoring the field.
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	
	// Decode the request body into the target destination.
	err := decoder.Decode(dst)
	if err != nil {
		// If the JSON contains syntax problem.
		// Syntax error can also invoke io.ErrUnexpectedEOF
		var syntaxError *json.SyntaxError
		// A JSON value is not appropriate for the destination Go type
		var unmarshalTypeError *json.UnmarshalTypeError
		// Decode destination is not valid (usually because it is not a pointer).
		// This is actually a problem with our application code, Not the JSON itself
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does then return a plain english error message
		// which contains the location of the problem
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Request contains badly formed JSON (at character %d)",
				syntaxError.Offset)

		// In Some cases Decode may also return an io.ErrUnexpectedEOF error for
		// syntax problem in the JSON. So we check for this using errors.IS() and
		// return a generic error message. There is an open issue reagarding this at
		// http://github.com/golang/go/issues/25956
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Re contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError erros. These occur when the
		// JSON value is th ewrong type for the target destination. If the error
		// relates to a specific field, then we include that in our error message to
		// make it easier for the client to debug
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field)
			} else {
				return fmt.Errorf("body contains incorrect JSON type (at character %d)",
					unmarshalTypeError.Offset)
			}
			
		// An io.EOF error will be returned by Decode() if the request body is empty
		case errors.Is(err, io.EOF):
			return errors.New("Request body must not be empty")
			
		// If the JSON contains a field which cannot be mapped to the target
		// destination then Decode() will now return an error message in the format
		// "json: unknown field '<name'". We check for this, extract the field name
		// from the error, and interpolate it into our custom error message. Not
		// that there's an open issue at https://github.com/golang/go/issues/29035
		// regarding turning this into a distinct type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body is too large (>1MB), Decode will fail.
		// NOTE: https://github.com/golang/go/issues/30715
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larget than %d bytes", maxbytes)
			
		// A json.InvalidUnmarshalError error will be returneed if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an
		// error to our handler.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	// Call decode again, using a pointer to an empty anonymous struct as the
	// destination to see if there is any other JSON value in the request body
	err = decoder.Decode(&struct{}{})
	
	// If err is not io.EOF then there are trailing value after that JSON
	// in the request body. But we don't want any value. So we return Error
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "Must be an valid integer value")
		return defaultValue
	}
	return i
}

func (app *application)background(fn func()) {
	app.wg.Add(1)
	// Launch a background goroutine
	go func () {
		defer app.wg.Done()
		
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		fn()
	}()
}
