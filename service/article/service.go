package article

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/store/article"
)

// Service implements app.ArticleService interface
// It serves articles from Postgres
type Service struct {
	store  app.ArticleStore
	router *chi.Mux
}

// New creates new Article service backed by Postgres
func New(store app.ArticleStore) *Service {
	router := chi.NewRouter()

	s := &Service{store: store, router: router}

	// Unauthenticated endpoints
	router.Get("/", s.HandleArticleList)
	router.Get("/{slug}", s.HandleArticleGet)

	return s
}

func NewFromDSN(DSN string) (*Service, error) {
	store, err := article.New(DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create article store for DSN %s", DSN)
	}

	return New(store), nil
}

func (s Service) List(f app.ArticleListFilter) ([]app.Article, error) {
	return s.store.List(f)
}

func (s Service) Get(slug string) (*app.Article, error) {
	return s.store.Get(slug)
}
