package article

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
	secret  []byte
}

func NewHTTP(store Store, secret []byte) (*Server, error) {
	s := &Server{
		router:  chi.NewRouter(),
		service: NewService(store),
		secret:  secret,
	}

	// Unauthenticated endpoints
	// s.router.Post("/", s.HandleUserRegister)
	// s.router.Post("/login", s.HandleUserLogin)
	// s.router.Get("/{slug}", server.HandleArticleGet)

	// Endpoints protected by JWT auth
	s.router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(s.secret))

		r.Post("/", s.HandleCreate)

		// r.Get("/", s.HandleUserGet)
		// r.Put("/", server.HandleUserUpdate)
		// r.Post("/{username}/follow", s.HandleProfileFollow)
		// r.Post("/{username}/unfollow", s.HandleProfileUnfollow)
	})

	return s, nil
}

// ServeHTTP implements http.handler interface and uses router ServeHTTP method
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ResponseSingle struct {
	Article app.Article `json:"article"`
}

func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		http.Error(w, transport.ServerError(app.ErrorUserNotInContext), http.StatusUnauthorized)
		return
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req CreateRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorUnmarshal, err), http.StatusUnprocessableEntity)
		return
	}

	a, err := s.service.Create(&req, &app.Profile{
		Name:  currentUser.Name,
		Bio:   currentUser.Bio,
		Image: currentUser.Image,
	})
	if err != nil {
		http.Error(w, transport.ServerError(err), http.StatusUnprocessableEntity)
		return
	}

	resp, err := json.Marshal(ResponseSingle{Article: *a})
	if err != nil {
		http.Error(w, transport.ServerError(transport.ErrorMarshal, err), http.StatusUnprocessableEntity)
		return
	}

	w.Write(resp)
}
