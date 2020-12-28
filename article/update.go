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
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
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
func (s *Service) Update(slug string, req *UpdateRequest) (*app.Article, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ValidationError(err)
	}

	// Find article
	a, err := s.store.GetArticle(slug)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get article for update"))
	}

	if a == nil {
		return nil, app.ServiceError(app.ErrorArticleNotExists)
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
	err = s.store.UpdateArticle(slug, a)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to update article"))
	}

	// Return updated article
	return s.store.GetArticle(slug)
}
