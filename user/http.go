package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

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
	s.router.Post("/", transport.WithError(s.HandleUserRegister))
	s.router.Post("/login", transport.WithError(s.HandleUserLogin))

	// Endpoints protected by JWT auth
	s.router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(s.secret, jwt.AuthTypeRequired))

		r.Get("/", transport.WithError(s.HandleUserGet))
		r.Put("/", transport.WithError(s.HandleUserUpdate))
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

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) error {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req LoginRequest
	err := decoder.Decode(&req)
	if err != nil {
		return app.ServiceError(errorInvalidRequest)
	}

	// Perform login in service
	user, err := s.service.Login(&req)
	if err != nil {
		return err
	}

	// Generate JWT
	token, err := jwt.New(user, s.secret)
	if err != nil {
		return app.InternalError(errors.Wrap(err, "jwt.New"))
	}

	// Prepare and send reply with user data, including token
	user.Token = token
	jsonUser, err := json.Marshal(Response{User: *user})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(jsonUser)
	return nil
}

func (s *Server) HandleUserRegister(w http.ResponseWriter, r *http.Request) error {
	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req RegisterRequest
	err := decoder.Decode(&req)
	if err != nil {
		return app.ServiceError(errorInvalidRequest)
	}

	// Perform register in service
	user, err := s.service.Register(&req)
	if err != nil {
		return err
	}

	// Generate JWT
	token, err := jwt.New(user, s.secret)
	if err != nil {
		return app.InternalError(errors.Wrap(err, "jwt.New"))
	}

	// Prepare and send reply with user data, including token
	user.Token = token
	jsonUser, err := json.Marshal(Response{User: *user})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonUser)
	return nil
}

// HandleUserGet gets the currently logged-in user. Requires authentication.
func (s *Server) HandleUserGet(w http.ResponseWriter, r *http.Request) error {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		return app.AuthError(app.ErrorUserNotInContext)
	}

	u, err := s.service.Get(currentUser.Email)
	if err != nil {
		return err
	}

	jsonUser, err := json.Marshal(Response{User: *u})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(jsonUser)
	return nil
}

// HandleUserUpdate changes currently logged-in user. Requires authentication.
func (s *Server) HandleUserUpdate(w http.ResponseWriter, r *http.Request) error {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		return app.AuthError(app.ErrorUserNotInContext)
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req UpdateRequest
	err := decoder.Decode(&req)
	if err != nil {
		return app.ServiceError(errorInvalidRequest)
	}

	// Check that user updates itself
	if req.User.Email != "" && req.User.Email != currentUser.Email {
		return app.AuthError(errorUserUpdateForbidden)
	}

	u, err := s.service.Update(currentUser.Id, &req)
	if err != nil {
		return err
	}

	// Regenerate JWT
	token, err := jwt.New(u, s.secret)
	if err != nil {
		return app.InternalError(errors.Wrap(err, "jwt.New"))
	}
	u.Token = token

	jsonUser, err := json.Marshal(Response{User: *u})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(jsonUser)
	return nil
}
