package article

import "github.com/dzeban/conduit/app"

func (s *Service) Get(slug string) (*app.Article, error) {
	a, err := s.store.GetArticle(slug)

	// Service store will return (nil, nil) when article not found.
	// Here, we set application level error to avoid nil dereference.
	if a == nil {
		return nil, app.ErrorArticleNotFound
	}

	return a, err
}
