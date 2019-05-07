package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
	"github.com/dzeban/conduit/postgres"
)

var (
	ErrJWTNoAuthorizationHeader = errors.New("no Authorization header")
	ErrJWTNoSignedClaim         = errors.New("token does not have signed claim")
	ErrJWTNoSubClaim            = errors.New("no sub claim")
)

// Server holds app server state
type Server struct {
	secret     []byte
	httpServer *http.Server
	articles   app.ArticlesService
	users      app.UserService
}

type ServerConfig struct {
	Port int `default:"8080"`
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
		http.Error(w, ServerError(err, "failed to list articles"), http.StatusInternalServerError)
		return
	}

	jsonArticles, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal json for articles list"), http.StatusInternalServerError)
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
		http.Error(w, ServerError(err, "failed to get article"), http.StatusInternalServerError)
		return
	}
	if article == nil {
		http.Error(w, ServerError(nil, fmt.Sprintf("article with slug %s not found", slug)), http.StatusNotFound)
		return
	}

	jsonArticle, err := json.Marshal(article)
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal json for article get"), http.StatusInternalServerError)
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
		http.Error(w, ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForRegister()
	if err != nil {
		http.Error(w, ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.users.Register(req)
	if err == app.ErrUserExists {
		http.Error(w, ServerError(err, "failed to register user"), http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to register user"), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
		"signed": true,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		http.Error(w, ServerError(err, "failed to create token"), http.StatusInternalServerError)
		return
	}

	user.Token = tokenString

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUser)
}

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req app.UserRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, ServerError(err, "failed to decode request"), http.StatusBadRequest)
		return
	}

	err = req.User.ValidateForLogin()
	if err != nil {
		http.Error(w, ServerError(err, "failed to validate request"), http.StatusBadRequest)
		return
	}

	user, err := s.users.Login(req)
	if err != nil {
		http.Error(w, ServerError(err, "failed to login"), http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.Email,
		"signed": true,
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		http.Error(w, ServerError(err, "failed to create token"), http.StatusInternalServerError)
		return
	}

	user.Token = tokenString

	// Prepare and send reply with user data, including token
	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonUser)
}

// HandleUserGet gets the currently logged-in user. Requires authentication.
func (s *Server) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value("email")
	if val == nil {
		http.Error(w, ServerError(nil, "no email in context"), http.StatusUnauthorized)
		return
	}

	email, ok := val.(string)
	if !ok {
		http.Error(w, ServerError(nil, "invalid auth email"), http.StatusUnauthorized)
		return
	}

	user, err := s.users.Get(email)
	if err == app.ErrUserNotFound {
		http.Error(w, ServerError(err, "failed to get user"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, ServerError(err, "failed to get user"), http.StatusInternalServerError)
		return
	}

	jsonUser, err := json.Marshal(app.UserRequest{User: *user})
	if err != nil {
		http.Error(w, ServerError(err, "failed to marshal response"), http.StatusInternalServerError)
		return
	}

	w.Write(jsonUser)
}

// jwtAuthHandler is a middleware that wraps next handler func with JWT token
// parsing and validation. It also stores authenticated user email into the
// context.
func (s *Server) jwtAuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header["Authorization"]
		if !ok {
			http.Error(w, ServerError(ErrJWTNoAuthorizationHeader, ""), http.StatusUnauthorized)
			return
		}

		claims, err := parseJWTClaimsFromHeader(authHeader[0], s.secret)
		if err != nil {
			http.Error(w, ServerError(err, "failed to parse jwt"), http.StatusBadRequest)
			return
		}

		if claims["signed"] != true {
			http.Error(w, ServerError(ErrJWTNoSignedClaim, ""), http.StatusUnauthorized)
			return
		}

		var sub interface{}
		if sub, ok = claims["sub"]; !ok {
			http.Error(w, ServerError(ErrJWTNoSubClaim, ""), http.StatusUnauthorized)
			return
		}

		// Store auth subject (email) to the context
		authCtx := context.WithValue(r.Context(), "email", sub)

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

func ServerError(err error, msg string) string {
	// errors.Wrap doesn't handle nil errors. To avoid nil pointer in error
	// message we create empty error here when error is nil
	if err == nil {
		err = errors.New("")
	}
	return fmt.Sprintf(`{"error":{"message":["%s"]}}`, errors.Wrap(err, msg))
}
