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
	"golang.org/x/crypto/bcrypt"
)

// Server holds app server state
type Server struct {
	httpServer *http.Server
	articles   app.ArticlesService
	users      app.UserService
}

type ServerConfig struct {
	Port int
}

// Config represents app configuration
type Config struct {
	Server   ServerConfig
	Articles app.ArticlesConfig
	Users    app.UsersConfig
}

// NewServer creates new server using config
func NewServer(conf Config) (*Server, error) {
	var err error

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: router,
	}

	var (
		articlesService app.ArticlesService
		userService     app.UserService
	)

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

	// Create articles service depending on configuration
	switch conf.Users.Type {
	case "mock":
		userService = mock.NewUserService()

	default:
		return nil, errors.Errorf("unknown users type '%s'", conf.Users.Type)
	}

	s := &Server{
		httpServer: httpServer,
		articles:   articlesService,
		users:      userService,
	}

	// Setup API endpoints
	router.HandleFunc("/articles/", s.HandleArticles).Methods("GET")
	router.HandleFunc("/articles/{slug}", s.HandleArticle).Methods("GET")

	router.HandleFunc("/users", s.HandleUserRegister).Methods("POST")

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

func (s *Server) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	type request struct {
		User app.User `json:user`
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req request
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	// Generate bcrypt hash from password
	user := req.User
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	// Replace plaintext password with bcrypt hashed
	user.Password = string(hashedPass)

	regUser, err := s.users.Register(user)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	// TODO: Generate token

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(regUser)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	w.WriteHeader(201) // 201: CREATED
	w.Write(jsonUser)
}
