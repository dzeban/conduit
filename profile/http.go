package profile

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
	router  *chi.Mux
	service *Service
}

func NewHTTP(store Store, secret []byte) (*Server, error) {
	s := &Server{
		router:  chi.NewRouter(),
		service: NewService(store),
	}

	s.router.
		With(jwt.Auth(secret, jwt.AuthTypeOptional)).
		Get("/{username}", transport.WithError(s.HandleGet))

	// Endpoints protected by JWT auth
	s.router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(secret, jwt.AuthTypeRequired))

		r.Post("/{username}/follow", transport.WithError(s.HandleFollow))
		r.Delete("/{username}/follow", transport.WithError(s.HandleUnfollow))
	})

	return s, nil
}

type Response struct {
	Profile app.Profile `json:"profile"`
}

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) error {
	username := chi.URLParam(r, "username")

	currentUser, _ := app.UserFromContext(r.Context())

	p, err := s.service.Get(username, currentUser)
	if err != nil {
		return err
	}

	jsonProfile, err := json.Marshal(Response{Profile: *p})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(jsonProfile)
	return nil
}

func (s *Server) HandleFollow(w http.ResponseWriter, r *http.Request) error {
	username := chi.URLParam(r, "username")

	currentUser, _ := app.UserFromContext(r.Context())

	p, err := s.service.Follow(currentUser, username)
	if err != nil {
		return err
	}

	jsonProfile, err := json.Marshal(Response{Profile: *p})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(jsonProfile)
	return nil
}

func (s *Server) HandleUnfollow(w http.ResponseWriter, r *http.Request) error {
	username := chi.URLParam(r, "username")

	currentUser, _ := app.UserFromContext(r.Context())

	p, err := s.service.Unfollow(currentUser, username)
	if err != nil {
		return err
	}

	jsonProfile, err := json.Marshal(Response{Profile: *p})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(jsonProfile)
	return nil
}
