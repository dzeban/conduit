package article

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		req CreateRequest
		err error
	}{
		{
			CreateRequest{
				ArticleRequest{
					Title: "",
					Body:  "",
				},
			},
			app.ErrorValidationTitleIsRequired,
		},
		{
			CreateRequest{
				ArticleRequest{
					Title: "",
					Body:  "x",
				},
			},
			app.ErrorValidationTitleIsRequired,
		},
		{
			CreateRequest{
				ArticleRequest{
					Title: "x",
					Body:  "",
				},
			},
			app.ErrorValidationBodyIsRequired,
		},
		{
			CreateRequest{
				ArticleRequest{
					Title:       "",
					Body:        "",
					Description: "x",
				},
			},
			app.ErrorValidationTitleIsRequired,
		},
		{
			CreateRequest{
				ArticleRequest{
					Title: "title",
					Body:  "body",
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		err := tt.req.Validate()
		if err != tt.err {
			t.Errorf("Validate(%+v): invalid error, expected '%v', got '%v'", tt.req, tt.err, err)
		}
	}
}

func TestCreate(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		req     *CreateRequest
		errType app.ErrorType
		err     error
	}{
		{
			"EmptyValidation",
			&CreateRequest{},
			app.ErrorTypeValidation,
			app.ErrorValidationTitleIsRequired,
		},
		// NOTE: We should check the case when existing article is created but
		// we can't do so because article uniquiness is checked by slug which is
		// generated from random value via uniuri.
		//
		// {
		// 		"Existing",
		// 		...
		// },
		{
			"Valid",
			&CreateRequest{
				ArticleRequest{
					Title:       "new",
					Description: "new",
					Body:        "new",
				},
			},
			0,
			nil,
		},
	}

	s := NewService(mock.NewArticleStore(), mock.NewProfilesStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Create(tt.req, &mock.Author)
			if err != nil {
				// Check error
				var e app.Error
				// Unwrap service.Error
				if !errors.As(err, &e) {
					t.Errorf("Create(%v): invalid error: expected %T, got %T", tt.req, e, err)
					return
				}

				// Check error type
				if e.Type != tt.errType {
					t.Errorf("Create(%v): invalid error type: expected %v, got %v", tt.req, e.Type, tt.errType)
					return
				}

				// Check error value
				if tt.err != nil {
					if e.Err != tt.err {
						t.Errorf("Create(%v): invalid error value: expected %v, got %v", tt.req, e.Err, tt.err)
						return
					}

				}
			}
		})

	}
}
