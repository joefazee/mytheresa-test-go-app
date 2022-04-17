package main

import (
	"github.com/joefazee/mytheresa/pkg/data"
	"github.com/joefazee/mytheresa/pkg/validator"
	"net/http"
)

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) productsHandler(w http.ResponseWriter, req *http.Request) {

	var input struct {
		data.Filter
	}

	v := validator.New()
	qs := req.URL.Query()

	input.Category = app.readString(qs, "category", "")
	input.PriceLessThan = app.readInt(qs, "priceLessThan", 0, v)
	input.Filter.Page = app.readInt(qs, "page", 1, v)
	input.Filter.PageSize = app.readInt(qs, "page_size", 5, v)

	if data.ValidateFilter(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, req, v.Errors)
		return
	}

	products, metadata, err := app.products.GetAll(input.Filter)
	if err != nil {
		app.serverErrorResponse(w, req, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"data": products.Marshall(), "metdata": metadata}, nil)
}
