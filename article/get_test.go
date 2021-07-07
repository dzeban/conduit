package article

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		slug string
		err  error
		want *app.Article
	}{
		{
			"NotFound",
			"xxx",
			errorArticleNotFound,
			nil,
		},
		{
			"Valid",
			mock.ArticleValid.Slug,
			nil,
			&mock.ArticleValid,
		},
	}

	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := s.Get(tt.slug)
			if err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("Get(%v): invalid error, want %v, got %v\n", tt.slug, tt.err, err)
				}
				return
			}
			aa := testArticle(*a)
			if !aa.Equals(tt.want) {
				t.Errorf("Get(%v): invalid article, want %v, got %v\n", tt.slug, tt.want, a)
			}
		})
	}
}
