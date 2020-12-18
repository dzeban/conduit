package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
	"github.com/dzeban/conduit/transport"
)

type Server struct {
	router *chi.Mux

	service *Service
	secret  []byte
}

func NewHTTP(store Store, secret []byte) (*Server, error) {
	s := &Server{
		router:  chi.NewRouter(),
		service: NewService(store),
		secret:  secret,
	}

	// Unauthenticated endpoints
	s.router.Post("/", s.HandleUserRegister)
	s.router.Post("/login", s.HandleUserLogin)

	// Endpoints protected by JWT auth
	s.router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(s.secret))

		r.Get("/", s.HandleUserGet)
		r.Put("/", s.HandleUserUpdate)
	})

	return s, nil
}

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Response represent a json structure used
// in user serice responses.
type Response struct {
	User app.User `json:"user"`
}

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req LoginRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorUnmarshal), http.StatusUnprocessableEntity)
		return
	}

	// Perform login in service
	user, err := s.service.Login(&req)
	if err != nil {
		// Return correct error to the client
		var e app.Error
		if !errors.As(err, &e) {
			http.Error(w, transport.ServerError(fmt.Errorf("invalid error data: %v", err)), http.StatusUnprocessableEntity)
			return
		}

		if e.Type == app.ErrorTypeService || e.Type == app.ErrorTypeValidation {
			http.Error(w, transport.ServerError(app.ErrorLogin, err), http.StatusUnauthorized)
			return
		}

		http.Error(w, transport.ServerError(app.ErrorInternal, err), http.StatusUnprocessableEntity)
		return
	}

	// Generate JWT
	token, err := jwt.New(user, s.secret)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorJWT, err), http.StatusUnprocessableEntity)
		return
	}

	// Prepare and send reply with user data, including token
	user.Token = token
	jsonUser, err := json.Marshal(Response{User: *user})
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorMarshal, err), http.StatusUnprocessableEntity)
		return
	}

	w.Write(jsonUser)
}

func (s *Server) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req RegisterRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorUnmarshal, err), http.StatusUnprocessableEntity)
		return
	}

	// Perform register in service
	user, err := s.service.Register(&req)
	if err != nil {
		http.Error(w, transport.ServerError(app.ErrorRegister, err), http.StatusUnprocessableEntity)
		return
	}

	// Generate JWT
	token, err := jwt.New(user, s.secret)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorJWT, err), http.StatusUnprocessableEntity)
		return
	}

	// Prepare and send reply with user data, including token
	user.Token = token
	jsonUser, err := json.Marshal(Response{User: *user})
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorMarshal, err), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUser)
}

// HandleUserGet gets the currently logged-in user. Requires authentication.
func (s *Server) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, transport.ServerError(app.ErrorUserNotInContext), http.StatusUnauthorized)
		return
	}

	u, err := s.service.Get(currentUser.Email)
	if err != nil {
		http.Error(w, transport.ServerError(err), http.StatusUnprocessableEntity)
		return
	}

	jsonUser, err := json.Marshal(Response{User: *u})
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorMarshal, err), http.StatusUnprocessableEntity)
		return
	}

	w.Write(jsonUser)
}

// HandleUserUpdate changes currently logged-in user. Requires authentication.
func (s *Server) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, transport.ServerError(app.ErrorUserNotInContext), http.StatusUnauthorized)
		return
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req UpdateRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorUnmarshal, err), http.StatusUnprocessableEntity)
		return
	}

	// Check that user updates itself
	if req.User.Email != "" && currentUser.Email != req.User.Email {
		http.Error(w, transport.ServerError(app.ErrorUserUpdateForbidden), http.StatusUnauthorized)
		return
	}

	u, err := s.service.Update(currentUser.Email, &req)
	if err != nil {
		http.Error(w, transport.ServerError(app.ErrorUpdate), http.StatusUnprocessableEntity)
		return
	}

	// Regenerate JWT
	token, err := jwt.New(u, s.secret)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorJWT, err), http.StatusUnprocessableEntity)
		return
	}
	u.Token = token

	jsonUser, err := json.Marshal(Response{User: *u})
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorMarshal, err), http.StatusUnprocessableEntity)
		return
	}

	w.Write(jsonUser)
}
