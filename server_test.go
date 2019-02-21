package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dzeban/conduit/mock"
)

////////////////////////////////////////////////////////
// Create server with HTTPTEST.NewServer ONCE
// Write tests for articles in a single function with a table-driven tests
////////////////////////////////////////////////////////

func TestArticlesList(t *testing.T) {
	mockArticlesService, _ := mock.New()
	s := Server{
		articles: mockArticlesService,
	}

	req := httptest.NewRequest("GET", "/articles", nil)

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

	req := httptest.NewRequest("GET", "/articles/xxx", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.HandleArticle)

	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusNotFound {
		t.Errorf("invalid status code: expected %v got %v'", http.StatusNotFound, status)
	}
}

func TestArticle(t *testing.T) {
	mockArticlesService, _ := mock.New()
	s := Server{
		articles: mockArticlesService,
	}

	req := httptest.NewRequest("GET", "/articles/slug-1", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.HandleArticle)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("invalid status code: expected %v got %v'", http.StatusOK, status)
	}

	body := rr.Body.String()

	expected := "Description 2"
	if !strings.Contains(body, expected) {
		t.Errorf("invalid body: expected '%v', got '%v'", expected, body)
	}
}
