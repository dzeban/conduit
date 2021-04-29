package article

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
	"github.com/dzeban/conduit/transport"
)

type Server struct {
	router  *chi.Mux
	service *Service
	secret  []byte
}

func NewHTTP(store Store, profilesStore ProfilesStore, secret []byte) (*Server, error) {
	s := &Server{
		router:  chi.NewRouter(),
		service: NewService(store, profilesStore),
		secret:  secret,
	}

	// Unauthenticated endpoints
	s.router.Get("/", transport.WithError(s.HandleList))
	s.router.Get("/{slug}", transport.WithError(s.HandleGet))

	// Endpoints protected by JWT auth
	s.router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(s.secret, jwt.AuthTypeRequired))

		r.Post("/", transport.WithError(s.HandleCreate))
		r.Get("/feed", transport.WithError(s.HandleFeed))
		r.Put("/{slug}", transport.WithError(s.HandleUpdate))
		r.Delete("/{slug}", transport.WithError(s.HandleDelete))
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

type ResponseMulti struct {
	Articles []*app.Article `json:"articles"`
	Count    int            `json:"articlesCount"`
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) error {
	slug := chi.URLParam(r, "slug")

	a, err := s.service.Get(slug)
	if err != nil {
		return err
	}

	resp, err := json.Marshal(ResponseSingle{Article: *a})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(resp)
	return nil
}

func (s *Server) HandleFeed(w http.ResponseWriter, r *http.Request) error {
	// Construct filter from query params
	params := r.URL.Query()
	filter := app.NewArticleListFilter()

	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		return app.AuthError(app.ErrorUserNotInContext)
	}

	filter.CurrentUser = currentUser

	if limit := params.Get("limit"); limit != "" {
		l, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return app.ServiceError(errorArticleInvalidLimit)
		}
		filter.Limit = l
	}

	if offset := params.Get("offset"); offset != "" {
		o, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return app.ServiceError(errorArticleInvalidOffset)
		}
		filter.Offset = o
	}

	// Get the article list from service
	articles, err := s.service.List(&filter)
	if err != nil {
		return err
	}

	// Marshal response
	resp, err := json.Marshal(ResponseMulti{
		Articles: articles,
		Count:    len(articles),
	})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(resp)
	return nil
}

func (s *Server) HandleList(w http.ResponseWriter, r *http.Request) error {
	// Construct filter from query params
	params := r.URL.Query()
	filter := app.NewArticleListFilter()
	if author := params.Get("author"); author != "" {
		filter.Author = &app.Profile{
			Name: author,
		}
	}

	if limit := params.Get("limit"); limit != "" {
		l, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return app.ServiceError(errorArticleInvalidLimit)
		}
		filter.Limit = l
	}

	if offset := params.Get("offset"); offset != "" {
		o, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return app.ServiceError(errorArticleInvalidOffset)
		}
		filter.Offset = o
	}

	// Get the article list from service
	articles, err := s.service.List(&filter)
	if err != nil {
		return err
	}

	// Marshal response
	resp, err := json.Marshal(ResponseMulti{
		Articles: articles,
		Count:    len(articles),
	})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(resp)
	return nil
}

func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) error {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		return app.AuthError(app.ErrorUserNotInContext)
	}

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req CreateRequest
	err := decoder.Decode(&req)
	if err != nil {
		return app.ServiceError(errorInvalidRequest)
	}

	author := app.Profile{
		Id:   currentUser.Id,
		Name: currentUser.Name,
	}
	a, err := s.service.Create(&req, &author)
	if err != nil {
		return err
	}

	resp, err := json.Marshal(ResponseSingle{Article: *a})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(resp)
	return nil
}

func (s *Server) HandleUpdate(w http.ResponseWriter, r *http.Request) error {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		return app.AuthError(app.ErrorUserNotInContext)
	}

	slug := chi.URLParam(r, "slug")

	// Decode user request from JSON body
	decoder := json.NewDecoder(r.Body)
	var req UpdateRequest
	err := decoder.Decode(&req)
	if err != nil {
		return app.ServiceError(errorInvalidRequest)
	}

	author := app.Profile{
		Id:   currentUser.Id,
		Name: currentUser.Name,
	}
	a, err := s.service.Update(slug, &author, &req)
	if err != nil {
		return err
	}

	resp, err := json.Marshal(ResponseSingle{Article: *a})
	if err != nil {
		return app.InternalError(errors.Wrap(err, "json.Marshal"))
	}

	w.Write(resp)
	return nil
}

func (s *Server) HandleDelete(w http.ResponseWriter, r *http.Request) error {
	currentUser, ok := app.UserFromContext(r.Context())
	if !ok {
		return app.AuthError(app.ErrorUserNotInContext)
	}

	slug := chi.URLParam(r, "slug")

	author := app.Profile{
		Id:   currentUser.Id,
		Name: currentUser.Name,
	}
	err := s.service.Delete(slug, &author)
	if err != nil {
		return err
	}

	w.Write(nil)
	return nil
}
