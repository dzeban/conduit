package article

import (
	"github.com/dchest/uniuri"
	"github.com/go-chi/chi"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
	"github.com/dzeban/conduit/store/article"
)

// Service implements app.ArticleService interface
// It serves articles from Postgres
type Service struct {
	store  app.ArticleStore
	router *chi.Mux
	secret []byte
}

// New creates new Article service backed by Postgres
func New(store app.ArticleStore, secret string) *Service {
	router := chi.NewRouter()

	s := &Service{
		store:  store,
		router: router,
		secret: []byte(secret),
	}

	// Unauthenticated endpoints
	router.Get("/", s.HandleArticleList)

	// Endpoints protected by JWT auth
	router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(s.secret))

		r.Get("/feed", s.HandleArticleFeed)
		r.Post("/", s.HandleArticleCreate)
		r.Delete("/{slug}", s.HandleArticleDelete)
	})

	router.Get("/{slug}", s.HandleArticleGet)

	return s
}

func NewFromDSN(DSN string, secret string) (*Service, error) {
	store, err := article.New(DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create article store for DSN %s", DSN)
	}

	return New(store, secret), nil
}

func (s Service) List(f app.ArticleListFilter) ([]app.Article, error) {
	return s.store.List(f)
}

func (s Service) Feed(f app.ArticleListFilter) ([]app.Article, error) {
	return s.store.Feed(f)
}

func (s Service) Get(slug string) (*app.Article, error) {
	return s.store.Get(slug)
}

func (s Service) Create(a *app.Article) error {
	a.Slug = slug.Make(a.Title) + "-" + uniuri.NewLen(4)
	return s.store.Create(a)
}

func (s Service) Delete(slug string) error {
	return s.store.Delete(slug)
}
