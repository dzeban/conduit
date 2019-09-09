package article

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/dzeban/conduit/app"
)

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func parseArticleListFilter(r *http.Request) app.ArticleListFilter {
	// Create filter with default values
	f := app.NewArticleListFilter()

	// Fill filter struct from query params
	q := r.URL.Query()

	s := q.Get("limit")
	if s != "" {
		val, err := strconv.ParseUint(s, 10, 64)
		// Ignore parsing error but set only if parsed
		if err == nil {
			f.Limit = val
		}
	}

	s = q.Get("offset")
	if s != "" {
		val, err := strconv.ParseUint(s, 10, 64)
		// Ignore parsing error but set only if parsed
		if err == nil {
			f.Offset = val
		}
	}

	f.Username = q.Get("author")

	return f
}

// HandleArticles is a handler for /articles API endpoint
func (s *Service) HandleArticleList(w http.ResponseWriter, r *http.Request) {
	f := parseArticleListFilter(r)

	articles, err := s.List(f)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to list articles"), http.StatusInternalServerError)
		return
	}

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal json for articles list"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticles)
}

// HandleArticles is a handler for /articles API endpoint
func (s *Service) HandleArticleFeed(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("username")
	if val == nil {
		http.Error(w, app.ServerError(nil, "no username in context"), http.StatusUnauthorized)
		return
	}

	username, ok := val.(string)
	if !ok {
		http.Error(w, app.ServerError(nil, "invalid auth email"), http.StatusUnauthorized)
		return
	}

	if username == "" {
		http.Error(w, app.ServerError(nil, "empty auth email"), http.StatusUnauthorized)
		return
	}

	f := parseArticleListFilter(r)

	f.Username = username

	articles, err := s.Feed(f)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to get articles feed"), http.StatusInternalServerError)
		return
	}

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal json for articles feed"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticles)
}

// HandleArticle is a handler for /article/{slug} API endpoint
func (s *Service) HandleArticleGet(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := s.Get(slug)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to get article"), http.StatusInternalServerError)
		return
	}
	if article == nil {
		http.Error(w, app.ServerError(nil, fmt.Sprintf("article with slug %s not found", slug)), http.StatusNotFound)
		return
	}

	jsonArticle, err := json.Marshal(article)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal json for article get"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticle)
}

func (s *Service) HandleArticleCreate(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("username")
	if val == nil {
		http.Error(w, app.ServerError(nil, "no username in context"), http.StatusUnauthorized)
		return
	}

	username, ok := val.(string)
	if !ok {
		http.Error(w, app.ServerError(nil, "invalid auth email"), http.StatusUnauthorized)
		return
	}

	if username == "" {
		http.Error(w, app.ServerError(nil, "empty auth email"), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req app.ArticleCreateRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	// Create article
	req.Article.Author.Name = username
	err = s.Create(&req.Article)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to create article"), http.StatusInternalServerError)
	}

	// Return newly created article
	article, err := s.Get(req.Article.Slug)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to get article"), http.StatusInternalServerError)
		return
	}
	if article == nil {
		http.Error(w, app.ServerError(nil, fmt.Sprintf("article with slug %s not found", req.Article.Slug)), http.StatusNotFound)
		return
	}

	jsonArticle, err := json.Marshal(article)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to marshal json for article get"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticle)
}

func (s *Service) HandleArticleDelete(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("username")
	if val == nil {
		http.Error(w, app.ServerError(nil, "no username in context"), http.StatusUnauthorized)
		return
	}

	username, ok := val.(string)
	if !ok {
		http.Error(w, app.ServerError(nil, "invalid auth email"), http.StatusUnauthorized)
		return
	}

	if username == "" {
		http.Error(w, app.ServerError(nil, "empty auth email"), http.StatusUnauthorized)
		return
	}

	slug := chi.URLParam(r, "slug")

	// Check that article exists
	article, err := s.Get(slug)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to get article"), http.StatusInternalServerError)
		return
	}
	if article == nil {
		http.Error(w, app.ServerError(nil, fmt.Sprintf("article with slug %s not found", slug)), http.StatusNotFound)
		return
	}

	// Check that this article belongs to this user
	fmt.Println(article.Author.Name, username)
	if article.Author.Name != username {
		http.Error(w, app.ServerError(err, "not allowed to delete other user article"), http.StatusForbidden)
		return
	}

	err = s.Delete(slug)
	if err != nil {
		http.Error(w, app.ServerError(err, "failed to delete article"), http.StatusInternalServerError)
	}
}
