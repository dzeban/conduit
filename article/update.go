package article

import (
	"time"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

type UpdateRequest struct {
	Article UpdateArticle `json:"article"`
}

type UpdateArticle struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Body        string `json:"body,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	if empty.MatchString(r.Article.Title) &&
		empty.MatchString(r.Article.Description) &&
		empty.MatchString(r.Article.Body) {
		return errors.New("at least one of title, description, body is required for update")
	}

	return nil
}

// Update modifies article found by slug with the new data in req.
// Returns updated article.
func (s *Service) Update(slug string, author *app.Profile, req *UpdateRequest) (*app.Article, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ServiceError(err)
	}

	// Find article
	a, err := s.store.GetArticle(slug)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get article for update"))
	}

	if a == nil {
		return nil, app.ServiceError(errorArticleNotFound)
	}

	// Check that article belongs to author
	if a.Author.Id != author.Id {
		return nil, app.ServiceError(errorArticleUpdateForbidden)
	}

	// Fill updated fields
	if !empty.MatchString(req.Article.Title) {
		a.Title = req.Article.Title
	}

	if !empty.MatchString(req.Article.Description) {
		a.Description = req.Article.Description
	}

	if !empty.MatchString(req.Article.Body) {
		a.Body = req.Article.Body
	}

	// Refresh updated timestamp
	a.Updated = time.Now()

	// Persist updated article in the store
	err = s.store.UpdateArticle(a)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to update article"))
	}

	// Return updated article
	a, err = s.store.GetArticle(slug)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get article after update"))
	}

	return a, nil
}
