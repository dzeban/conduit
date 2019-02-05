package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Server holds app server state
type Server struct {
	db         *sqlx.DB
	httpServer *http.Server
}

// Config represents app configuration
type Config struct {
	Port int
	DSN  string
}

// NewServer creates new server using config
func NewServer(conf Config) *Server {
	// Open database connection and store it into server
	db, err := sqlx.Open("postgres", conf.DSN)
	if err != nil {
		panic("can't connect to db")
	}

	// Create a dedicated HTTP server to set listen address
	// gin default router is used as the handler
	router := gin.Default()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: router,
	}

	s := &Server{
		db:         db,
		httpServer: httpServer,
	}

	// Setup API endpoints
	router.GET("/articles/", s.HandleArticles)
	router.GET("/articles/:slug", s.HandleArticle)

	return s
}

// Run starts server to listen and serve requests
func (s *Server) Run() {
	s.httpServer.ListenAndServe()
}

// HandleArticles is a handler for /articles API endpoint
func (s *Server) HandleArticles(c *gin.Context) {
	rows, err := s.db.Queryx(queryArticles)
	if err != nil {
		c.Status(500)
		return
	}

	var articles Articles
	for rows.Next() {
		var article Article
		err = rows.StructScan(&article)
		if err != nil {
			log.Println(err)
			continue
		}

		articles = append(articles, article)
	}

	c.JSON(200, articles)
}

// HandleArticle is a handler for /article/:slug API endpoint
func (s *Server) HandleArticle(c *gin.Context) {
	row := s.db.QueryRowx(queryArticle, c.Param("slug"))

	var article Article
	err := row.StructScan(&article)
	if err == sql.ErrNoRows {
		c.Status(404)
		return
	} else if err != nil {
		c.Status(500)
		return
	}

	c.JSON(200, article)
}
