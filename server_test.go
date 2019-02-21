package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dzeban/conduit/mock"
)

func TestArticlesList(t *testing.T) {
	mockArticlesService, _ := mock.New()
	s := Server{
		articles: mockArticlesService,
	}

	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.HandleArticles)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("invalid status code: expected %v got %v'", http.StatusOK, status)
	}

	body := rr.Body.String()

	expected := "Title 1"
	if !strings.Contains(body, expected) {
		t.Errorf("invalid body: expected %v, got %v", expected, body)
	}
}

func TestNotFoundArticle(t *testing.T) {
	mockArticlesService, _ := mock.New()
	s := Server{
		articles: mockArticlesService,
	}

	req, err := http.NewRequest("GET", "/article/xxx", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.HandleArticle)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusNotFound {
		t.Errorf("invalid status code: expected %v got %v'", http.StatusNotFound, status)
	}
}
