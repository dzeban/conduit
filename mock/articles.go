package mock

import (
	"time"

	"github.com/dzeban/conduit/app"
)

// ArticleService implements app.ArticleService interface
// It serves articles from memory
type ArticlesService struct {
	articles []app.Article
}

// New creates new Articles service
func New() (*ArticlesService, error) {
	return &ArticlesService{
		articles: []app.Article{
			{
				Title:       "Title 1",
				Description: "Description 1",
				Slug:        "slug-1",
				Body:        "Body 1",
				Created:     time.Now(),
				Updated:     time.Now(),
			},
			{
				Title:       "Title 2",
				Description: "Description 2",
				Slug:        "slug-2",
				Body:        "Body 2",
				Created:     time.Now(),
				Updated:     time.Now(),
			},
		},
	}, nil
}

// List returns n articles
func (s *ArticlesService) List(n int) ([]app.Article, error) {
	return s.articles, nil
}

// Get returns a single article by its slug
func (s *ArticlesService) Get(slug string) (*app.Article, error) {
	for _, a := range s.articles {
		if a.Slug == slug {
			return &a, nil
		}
	}
	return nil, nil
}
