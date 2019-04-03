package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
	"github.com/dzeban/conduit/postgres"
)

// Server holds app server state
type Server struct {
	secret     []byte
	httpServer *http.Server
	articles   app.ArticlesService
	users      app.UserService
}

type ServerConfig struct {
	Port int
}

// Config represents app configuration
type Config struct {
	Secret   string
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
		secret:     []byte(conf.Secret),
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
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	user := req.User

	err = user.ValidateForRegister()
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	err = s.users.Register(user, user.Password)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"signed": true,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	user.Token = tokenString

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: user})
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	w.WriteHeader(201) // 201: CREATED
	w.Write(jsonUser)
}
