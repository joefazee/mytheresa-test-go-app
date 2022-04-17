package main

import (
	"github.com/joefazee/mytheresa/pkg/logger"
	"github.com/joefazee/mytheresa/pkg/models/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	return &testServer{ts}
}

func newTestApp(t *testing.T) *application {

	return &application{
		config:   config{},
		products: mock.ProductModel{},
		logger:   logger.New(io.Discard, logger.LevelOff),
	}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body

}
