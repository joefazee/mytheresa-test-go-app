package main

import "net/http"

// routes() register all http routes and return http.Handler
func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/products", app.productsHandler)

	return mux
}
