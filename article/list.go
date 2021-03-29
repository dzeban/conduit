package article

import (
	"github.com/dzeban/conduit/app"
	"github.com/pkg/errors"
)

func (s *Service) List(filter *app.ArticleListFilter) ([]*app.Article, error) {
	// Validate filter
	err := filter.Validate()
	if err != nil {
		return nil, app.ValidationError(err)
	}

	// Fill author id in filter
	if filter.Author != nil {
		author, err := s.profileStore.GetProfile(filter.Author.Name, app.ProfileFromUser(filter.CurrentUser))
		if err != nil {
			return nil, app.InternalError(errors.Wrap(err, "failed to get author profile"))
		}

		filter.Author.Id = author.Id
	}

	return s.store.ListArticles(filter)
}
