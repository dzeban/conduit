package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	case "postgres":
		userService, err = postgres.NewUserService(conf.Users.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create users service")
		}

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
	router.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
	router.HandleFunc("/users", s.jwtAuthHandler(s.HandleUserGet)).Methods("GET")

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

	err = req.User.ValidateForRegister()
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	user, err := s.users.Register(req)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
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
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	w.WriteHeader(201) // 201: CREATED
	w.Write(jsonUser)
}

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	err = req.User.ValidateForLogin()
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	user, err := s.users.Login(req)
	if err != nil {
		w.WriteHeader(401)
		log.Println(err)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
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
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonUser)
}

// HandleUserGet gets the currently logged-in user. Requires authentication.
func (s *Server) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("email")
	if val == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	email, ok := val.(string)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, "invalid auth email")
		return
	}

	user, err := s.users.Get(email)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)
}

// jwtAuthHandler is a middleware that wraps next handler func with JWT token
// parsing and validation. It also stores authenticated user email into the
// context.
func (s *Server) jwtAuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := parseJWTClaimsFromHeader(authHeader[0], s.secret)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, err)
			return
		}

		if claims["signed"] != true {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, "token does not have signed claim")
			return
		}

		if claims["sub"] == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"error":{"body":["%s"]}}`, "no sub claim")
			return
		}

		// Store auth subject (email) to the context
		authCtx := context.WithValue(r.Context(), "email", claims["sub"])

		next(w, r.WithContext(authCtx))
	}
}

// parseJWTClaimsFromHeader takes JWT from Authorization header, parses and
// validates it and returns claims. JWT is expected in "Token <token>" format.
func parseJWTClaimsFromHeader(header string, secret []byte) (map[string]interface{}, error) {
	tokenVals := strings.Split(header, " ")

	if len(tokenVals) != 2 {
		return nil, errors.New("invalid auth header format, expected 2 elements")
	}

	if tokenVals[0] != "Token" {
		return nil, fmt.Errorf("invalid auth header format, expected Token <token>, got %#v", header)
	}

	token, err := jwt.Parse(tokenVals[1], func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt parsing error")
	}

	if !token.Valid {
		return nil, errors.New("jwt is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to cast claims to jwt.MapClaims")
	}

	return claims, nil
}
