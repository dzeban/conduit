package article

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// HandleArticles is a handler for /articles API endpoint
func (s *Service) HandleArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := s.List(20)
	if err != nil {
		// http.Error(w, ServerError(err, "failed to list articles"), http.StatusInternalServerError)
		http.Error(w, "failed to list articles", http.StatusInternalServerError)
		return
	}

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		// http.Error(w, ServerError(err, "failed to marshal json for articles list"), http.StatusInternalServerError)
		http.Error(w, "failed to marshal json for articles list", http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticles)
}

// HandleArticle is a handler for /article/{slug} API endpoint
func (s *Service) HandleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	article, err := s.Get(slug)
	if err != nil {
		// http.Error(w, ServerError(err, "failed to get article"), http.StatusInternalServerError)
		http.Error(w, "failed to get article", http.StatusInternalServerError)
		return
	}
	if article == nil {
		// http.Error(w, ServerError(nil, fmt.Sprintf("article with slug %s not found", slug)), http.StatusNotFound)
		http.Error(w, fmt.Sprintf("article with slug %s not found", slug), http.StatusNotFound)
		return
	}

	jsonArticle, err := json.Marshal(article)
	if err != nil {
		// http.Error(w, ServerError(err, "failed to marshal json for article get"), http.StatusInternalServerError)
		http.Error(w, "failed to marshal json for article get", http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticle)
}
