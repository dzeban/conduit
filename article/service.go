package article

import (
	"regexp"

	"github.com/dzeban/conduit/app"
)

// ArticleStore defines an interface to work with articles
type Store interface {
	CreateArticle(a *app.Article) error
	GetArticle(slug string) (*app.Article, error)
	ListArticles(f *app.ArticleListFilter) ([]*app.Article, error)
	UpdateArticle(a *app.Article) error
	DeleteArticle(id int) error
}

// ProfilesStore provides helper to get author with all its fields (like id) by
// username
type ProfilesStore interface {
	GetProfile(username string, follower *app.User) (*app.Profile, error)
}

// Service provides methods for articles
type Service struct {
	store        Store
	profileStore ProfilesStore
}

// NewService creates new instance of the service with provided store
func NewService(store Store, profileStore ProfilesStore) *Service {
	return &Service{store, profileStore}
}

// empty is regexp to validate for "empty" string.
// Empty string is the one with zero length or containing only whitespaces.
var empty = regexp.MustCompile(`^[[:space:]]*$`)
