package article

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		slug    string
		req     UpdateRequest
		errType app.ErrorType
		err     error
		article *app.Article
	}{
		{
			"Validation",
			"xxx",
			UpdateRequest{
				UpdateArticle{
					Title:       " ",
					Description: "\t",
					Body:        "\n",
				},
			},
			app.ErrorTypeService,
			nil,
			nil,
		},
		{
			"NonExisting",
			"absent",
			UpdateRequest{
				UpdateArticle{
					Title: "new",
				},
			},
			app.ErrorTypeService,
			errorArticleNotFound,
			nil,
		},
		{
			"WhitespaceField",
			mock.ArticleUpdated.Slug,
			UpdateRequest{
				UpdateArticle{
					Title:       "Updated title",
					Description: "  ", // this should not trigger description update
				},
			},
			0,
			nil,
			&app.Article{
				Slug:        mock.ArticleUpdated.Slug,
				Title:       "Updated title",
				Description: mock.ArticleUpdated.Description,
				Body:        mock.ArticleUpdated.Body,
				Author:      mock.ArticleUpdated.Author,
			},
		},
		{
			"Valid",
			mock.ArticleValid.Slug,
			UpdateRequest{
				UpdateArticle{
					Title:       "New title",
					Description: "New description",
					Body:        "New body",
				},
			},
			0,
			nil,
			&app.Article{
				Slug:        mock.ArticleUpdated.Slug,
				Title:       "New title",
				Description: "New description",
				Body:        "New body",
				Author:      mock.ArticleUpdated.Author,
			},
		},
	}

	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := s.Update(tt.slug, &mock.Profile1, &tt.req)
			if err != nil {
				// Check error
				var e app.Error

				// Unwrap service.Error
				if !errors.As(err, &e) {
					t.Errorf("Update(%v): invalid error: expected %T, got %T", tt.req, e, err)
					return
				}

				// Check error type
				if e.Type != tt.errType {
					t.Errorf("Update(%v): invalid error type: expected %v, got %v", tt.req, tt.errType, e.Type)
					return
				}

				// Check error value
				if tt.err != nil {
					if e.Err != tt.err {
						t.Errorf("Update(%v): invalid error value: expected %v, got %v", tt.req, tt.err, e.Err)
						return
					}
				}
			}

			if tt.article != nil {
				// Check returned article
				aa := testArticle(*a)
				if !aa.Equals(tt.article) {
					t.Errorf("response not matched, expected '%+v', got '%+v'", tt.article, a)
				}
			}
		})
	}
}

func TestUpdateTimestamp(t *testing.T) {
	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	prevUpdated := mock.ArticleValid.Updated
	a, err := s.Update(mock.ArticleValid.Slug, &mock.Profile1, &UpdateRequest{
		UpdateArticle{
			Title: "new title",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !a.Updated.After(prevUpdated) {
		t.Errorf("Updated timestamp was not set after Update(): prev %v, current %v", prevUpdated, a.Updated)
	}
}

func TestUpdateAuthorCheck(t *testing.T) {
	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	invalidAuthor := app.Profile{
		Id:   999,
		Name: "Evil",
	}
	_, err := s.Update(mock.ArticleValid.Slug, &invalidAuthor, &UpdateRequest{
		UpdateArticle{
			Title: "new title",
		},
	})

	if !errors.Is(err, errorArticleUpdateForbidden) {
		t.Errorf("invalid error, expected '%v', got '%v'", errorArticleUpdateForbidden, err)
	}
}
