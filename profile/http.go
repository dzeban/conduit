package profile

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

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
		Get("/{username}", s.HandleGet)

	// Endpoints protected by JWT auth
	s.router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(secret, jwt.AuthTypeRequired))

		r.Post("/{username}/follow", s.HandleFollow)
		r.Delete("/{username}/unfollow", s.HandleUnfollow)
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

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	currentUser, _ := app.UserFromContext(r.Context())

	p, err := s.service.Get(username, currentUser)
	if err != nil {
		http.Error(w, transport.ServerError(err), http.StatusUnprocessableEntity)
		return
	}

	jsonProfile, err := json.Marshal(Response{Profile: *p})
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorMarshal, err), http.StatusUnprocessableEntity)
		return
	}

	w.Write(jsonProfile)
}

func (s *Server) HandleFollow(w http.ResponseWriter, r *http.Request)   {}
func (s *Server) HandleUnfollow(w http.ResponseWriter, r *http.Request) {}
