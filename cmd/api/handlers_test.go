package main

import (
	"encoding/json"
	"github.com/joefazee/mytheresa/pkg/data"
	"github.com/joefazee/mytheresa/pkg/models"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	healthcheck(rr, r)

	res := rr.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, res.StatusCode)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyStr := string(body)
	if bodyStr != "OK" {
		t.Errorf("want body to equal OK; got %q", bodyStr)
	}

}

func TestProducts(t *testing.T) {
	app := newTestApp(t)

	srv := newTestServer(t, app.routes())
	defer srv.Close()

	testTable := []struct {
		name          string
		urlPath       string
		wantCode      int
		wantBodyCount int
	}{
		{name: "no query string", urlPath: "/v1/products", wantCode: http.StatusOK, wantBodyCount: 5},
		{name: "get boots category", urlPath: "/v1/products?category=boots", wantCode: http.StatusOK, wantBodyCount: 3},
		{name: "price less than 9800", urlPath: "/v1/products?priceLessThan=89000", wantCode: http.StatusOK, wantBodyCount: 4},
		{name: "price less than 9800 AND category = ", urlPath: "/v1/products?priceLessThan=89000&category=sneakers", wantCode: http.StatusOK, wantBodyCount: 1},
		{name: "invalid category", urlPath: "/v1/products?category=sss", wantCode: http.StatusOK, wantBodyCount: 0},
		{name: "invalid priceLessThan data type", urlPath: "/v1/products?priceLessThan=sss", wantCode: http.StatusUnprocessableEntity, wantBodyCount: 0},
		{name: "must be greater than zero", urlPath: "/v1/products?page=0", wantCode: http.StatusUnprocessableEntity, wantBodyCount: 0},
		{name: "must not be more than 10m", urlPath: "/v1/products?page=10000001", wantCode: http.StatusUnprocessableEntity, wantBodyCount: 0},

		{name: "must be greater than zero", urlPath: "/v1/products?page_size=0", wantCode: http.StatusUnprocessableEntity, wantBodyCount: 0},
		{name: "must be a maximum of 100", urlPath: "/v1/products?page_size=101", wantCode: http.StatusUnprocessableEntity, wantBodyCount: 0},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			code, headers, body := srv.get(t, tt.urlPath)

			if headers.Get("Content-Type") != "application/json" {
				t.Errorf("Content-Type expected to be application/json; got %s", headers.Get("Content-Type"))
			}

			if code != tt.wantCode {
				t.Errorf("expected status code should be %d; got %d", tt.wantCode, code)
			}

			res := struct {
				Data    []models.PublicProduct
				Metdata data.Metadata
			}{}

			err := json.Unmarshal(body, &res)
			if err != nil {
				t.Errorf("%q", err)
			}

			if len(res.Data) != tt.wantBodyCount {
				t.Errorf("expected result count to be %d; got %d", tt.wantBodyCount, len(res.Data))
			}

		})
	}

}
