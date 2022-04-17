package main

import (
	"encoding/json"
	"github.com/joefazee/mytheresa/pkg/validator"
	"net/http"
	"net/url"
	"strconv"
)

type envelope map[string]interface{}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {

	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, mesage interface{}) {
	env := envelope{"error": mesage}

	if err := app.writeJSON(w, status, env, nil); err != nil {
		app.logger.Error(err, nil)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err, nil)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
