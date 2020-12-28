package article

import (
	"regexp"

	"github.com/dzeban/conduit/app"
)

// ArticleStore defines an interface to work with articles
type Store interface {
	ListArticles(f app.ArticleListFilter) ([]app.Article, error)
	GetArticle(slug string) (*app.Article, error)
	CreateArticle(a *app.Article) error
	UpdateArticle(slug string, a *app.Article) error
	DeleteArticle(slug string) error
}

// Service provides a service for interacting with user accounts
type Service struct {
	store Store
}

// NewService creates new instance of the service with provided store
func NewService(store Store) *Service {
	return &Service{store}
}

// empty is regexp to validate for "empty" string.
// Empty string is the one with zero length or containing only whitespaces.
var empty = regexp.MustCompile(`^[[:space:]]*$`)
