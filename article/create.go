package article

import (
	"time"

	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

const (
	// Length of random string appended to slug to make it unique
	SlugRandLen = 4
)

type CreateRequest struct {
	Article ArticleRequest `json:"article"`
}

type ArticleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	// TagList []Tag `json:"tagList"`
}

func (r *CreateRequest) Validate() error {
	if empty.MatchString(r.Article.Title) {
		return errorValidationTitleIsRequired
	}

	if empty.MatchString(r.Article.Body) {
		return errorValidationBodyIsRequired
	}

	return nil
}

// Create creates new article in the articles store
func (s *Service) Create(req *CreateRequest, author *app.Profile) (*app.Article, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ServiceError(err)
	}

	article := &app.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		Author:      *author,
		Created:     time.Now(),
		Updated:     time.Now(),

		// Generate slug
		Slug: slug.Make(req.Article.Title) + "-" + uniuri.NewLen(SlugRandLen),
	}

	// NOTE: we don't need to check if article exists because article is
	// identified by slug which is randomly generated

	// Persist article in the store
	err = s.store.CreateArticle(article)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to create article"))
	}

	return article, nil
}
