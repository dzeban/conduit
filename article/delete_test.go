package article

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name   string
		slug   string
		author *app.Profile
		err    error
	}{
		{
			"NotExisting",
			"xxx",
			&mock.Profile1,
			errorArticleNotFound,
		},
		{
			"Forbidden",
			mock.ArticleValid.Slug,
			&app.Profile{
				Id:   999,
				Name: "Evil",
			},
			errorArticleDeleteForbidden,
		},
		{
			"Valid",
			mock.ArticleValid.Slug,
			&mock.Profile1,
			nil,
		},
	}

	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Delete(tt.slug, tt.author)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("Delete(%v, %v): invalid error: expected '%v', got '%v'", tt.slug, tt.author, tt.err, err)
				}
			}
		})
	}
}

func TestDeleteForReal(t *testing.T) {
	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	err := s.Delete(mock.ArticleValid.Slug, &mock.Profile1)
	if err != nil {
		t.Errorf("Delete(%v, %v): unexpected error '%v'", mock.ArticleValid.Slug, mock.Profile1, err)
	}

	_, err = s.Get(mock.ArticleValid.Slug)
	if !errors.Is(err, errorArticleNotFound) {
		t.Errorf("Expected article not found after delete, got err '%v'", err)
	}
}
