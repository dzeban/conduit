package main

import (
	"fmt"
	"net/http"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/postgres"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Server holds app server state
type Server struct {
	httpServer *http.Server
	articles   app.ArticlesService
}

// Config represents app configuration
type Config struct {
	Port int
	DSN  string
}

// NewServer creates new server using config
func NewServer(conf Config) (*Server, error) {
	// Create a dedicated HTTP server to set listen address
	// gin default router is used as the handler
	router := gin.Default()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	// Create articles service depending on configuration
	// Currently, we have only 1 articles service based on Postgres
	// but later we may choose different implementation based on config
	articlesService, err := postgres.New(conf.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create articles service")
	}

	s := &Server{
		httpServer: httpServer,
		articles:   articlesService,
	}

	// Setup API endpoints
	router.GET("/articles/", s.HandleArticles)
	router.GET("/articles/:slug", s.HandleArticle)

	return s, nil
}

// Run starts server to listen and serve requests
func (s *Server) Run() {
	s.httpServer.ListenAndServe()
}

// HandleArticles is a handler for /articles API endpoint
func (s *Server) HandleArticles(c *gin.Context) {
	articles, err := s.articles.List(20)
	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, articles)
}

// HandleArticle is a handler for /article/:slug API endpoint
func (s *Server) HandleArticle(c *gin.Context) {
	article, err := s.articles.Get(c.Param("slug"))
	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, article)
}
