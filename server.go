package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
	"github.com/dzeban/conduit/postgres"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Server holds app server state
type Server struct {
	httpServer *http.Server
	articles   app.ArticlesService
}

type ServerConfig struct {
	Port int
}

// Config represents app configuration
type Config struct {
	Server   ServerConfig
	Articles app.ArticlesConfig
}

// NewServer creates new server using config
func NewServer(conf Config) (*Server, error) {
	var err error

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: router,
	}

	var articlesService app.ArticlesService

	// Create articles service depending on configuration
	switch conf.Articles.Type {
	case "postgres":
		articlesService, err = postgres.New(conf.Articles.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create articles service")
		}

	case "mock":
		articlesService, _ = mock.New()

	default:
		return nil, errors.Errorf("unknown article type '%s'", conf.Articles.Type)
	}

	s := &Server{
		httpServer: httpServer,
		articles:   articlesService,
	}

	// Setup API endpoints
	router.HandleFunc("/articles/", s.HandleArticles).Methods("GET")
	router.HandleFunc("/articles/{slug}", s.HandleArticle).Methods("GET")

	return s, nil
}

// Run starts server to listen and serve requests
func (s *Server) Run() {
	s.httpServer.ListenAndServe()
}

// HandleArticles is a handler for /articles API endpoint
func (s *Server) HandleArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := s.articles.List(20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticles)
}

// HandleArticle is a handler for /article/{slug} API endpoint
func (s *Server) HandleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	article, err := s.articles.Get(slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if article == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "article with slug %s not found", slug)
		return
	}

	jsonArticle, err := json.Marshal(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonArticle)
}
