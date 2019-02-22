package main

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dzeban/conduit/app"
)

var server *Server

func TestMain(m *testing.M) {
	mockConf := Config{
		Articles: app.ArticlesConfig{
			Type: "mock",
		},
	}

	var err error
	server, err = NewServer(mockConf)
	if err != nil {
		panic("failed to create server")
	}

	os.Exit(m.Run())
}

func TestArticles(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		status int
		body   string
	}{
		{"list", "GET", "/articles/", 200, "Title 1"},
		{"notfound", "GET", "/articles/xxx", 404, ""},
		{"single", "GET", "/articles/slug-2", 200, "Description 2"},
	}

	for _, expected := range tests {
		// Run in a subtest to distinguish tests by name
		t.Run(expected.name, func(t *testing.T) {
			req := httptest.NewRequest(expected.method, expected.url, nil)

			rr := httptest.NewRecorder()
			server.httpServer.Handler.ServeHTTP(rr, req)

			// Check status
			status := rr.Code
			if status != expected.status {
				t.Errorf("invalid status code: expected %v got %v'", expected.status, status)
			}

			// Check body
			body := rr.Body.String()
			if !strings.Contains(body, expected.body) {
				t.Errorf("invalid body: expected %v, got %v", expected.body, body)
			}
		})
	}
}
