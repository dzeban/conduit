package article

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

func (s *Service) Delete(slug string, author *app.Profile) error {
	// Find article to get its id and check author
	a, err := s.store.GetArticle(slug)
	if err != nil {
		return app.InternalError(errors.Wrap(err, "failed to get article for delete"))
	}

	if a == nil {
		return app.ServiceError(errorArticleNotFound)
	}

	// Check that article belongs to author
	if a.Author.Id != author.Id {
		return app.ServiceError(errorArticleDeleteForbidden)
	}

	err = s.store.DeleteArticle(a.Id)
	if err != nil {
		return app.InternalError(errors.Wrap(err, "article delete failed"))
	}

	return nil
}
