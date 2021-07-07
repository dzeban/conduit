package article

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
	"github.com/go-test/deep"
)

func TestListFilterValidation(t *testing.T) {
	tests := []struct {
		name    string
		filter  *app.ArticleListFilter
		errType app.ErrorType
	}{
		{
			"Limit0",
			&app.ArticleListFilter{CurrentUser: nil, Author: nil, Limit: 0, Offset: 1},
			app.ErrorTypeService,
		},
		{
			"LimitExceed",
			&app.ArticleListFilter{CurrentUser: nil, Author: nil, Limit: 9999999, Offset: 1},
			app.ErrorTypeService,
		},
		{
			"OffsetExceed",
			&app.ArticleListFilter{CurrentUser: nil, Author: nil, Limit: 1, Offset: 99999999},
			app.ErrorTypeService,
		},
		{
			"Valid",
			&app.ArticleListFilter{CurrentUser: nil, Author: nil, Limit: 1, Offset: 0},
			-1,
		},
	}

	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.List(tt.filter)
			if err != nil {
				var e app.Error
				if !errors.As(err, &e) {
					t.Errorf("List(%v): invalid error: expected %T, got %T", tt.filter, e, err)
				}

				if e.Type != tt.errType {
					t.Errorf("List(%v): invalid error type: expected %v, got %v", tt.filter, tt.errType, e.Type)
				}
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		filter   *app.ArticleListFilter
		expected []*app.Article
	}{
		{
			"Public",
			&app.ArticleListFilter{
				CurrentUser: nil,
				Author:      nil,
				Limit:       100,
				Offset:      0,
			},
			[]*app.Article{&mock.ArticleValid, &mock.ArticleUpdated, &mock.Article3},
		},
		{
			"ByAuthor",
			&app.ArticleListFilter{
				CurrentUser: nil,
				Author:      &mock.Profile2,
				Limit:       100,
				Offset:      0,
			},
			[]*app.Article{&mock.Article3},
		},
	}
	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as, err := s.List(tt.filter)
			if err != nil {
				t.Errorf("List(%v): unexpected error %v\n", tt.filter, err)
			}

			if diff := deep.Equal(as, tt.expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}
