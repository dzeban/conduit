package article

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
	"github.com/dzeban/conduit/mock"
)

const testSecret = "test"

// Equals implements custom comparison of articles needed in tests.
// It compares all of the fields except slug because it's generated with
// random and created/updated timestamps because they are not fixed.
type testArticle app.Article

func (a *testArticle) Equals(b *app.Article) bool {
	if a.Title != b.Title {
		return false
	}
	if a.Description != b.Description {
		return false
	}
	if a.Body != b.Body {
		return false
	}
	if a.Author.Name != b.Author.Name {
		return false
	}
	if a.Author.Bio != b.Author.Bio {
		return false
	}
	if a.Author.Image != b.Author.Image {
		return false
	}

	return true
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name        string
		contextUser *app.User
		body        string
		status      int
		resp        *ResponseSingle
		err         error
	}{
		{
			"Null",
			&app.User{},
			"",
			http.StatusUnauthorized,
			nil,
			// Don't check concrete error because it is returned by JWT middleware.
			// We need to be sure that unauthorized request returns error.
			nil,
		},
		{
			"Unmarshal",
			&mock.UserValid,
			"{",
			http.StatusUnprocessableEntity,
			nil,
			errorInvalidRequest,
		},
		{
			"InvalidRequest",
			&mock.UserValid,
			`{"article":{"title":"","description":"some","body":"some"}}`,
			http.StatusUnprocessableEntity,
			nil,
			errorValidationTitleIsRequired,
		},
		{
			"Valid",
			&mock.UserValid,
			`{"article":{"title":"new","description":"new","body":"new"}}`,
			http.StatusOK,
			&ResponseSingle{
				app.Article{
					Title:       "new",
					Description: "new",
					Body:        "new",
					Author: app.Profile{
						Name:  mock.UserValid.Name,
						Bio:   mock.UserValid.Bio,
						Image: mock.UserValid.Image,
					},
				},
			},
			nil,
		},
	}

	s, err := NewHTTP(mock.NewArticleStore(), mock.NewProfilesStore(), []byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwt.New(tt.contextUser, []byte(testSecret))
			if err != nil {
				t.Errorf("failed to make JWT")
				return
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Add("Authorization", "Token "+token)

			s.ServeHTTP(rr, req)

			resp := rr.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != tt.status {
				t.Errorf("incorrect status, expected %v, got %v", tt.status, resp.StatusCode)
				t.Errorf("resp body: %v", string(body))
				return
			}

			if tt.err != nil {
				if !strings.Contains(string(body), tt.err.Error()) {
					t.Errorf("expected error not found, expected '%v', got %s", tt.err.Error(), body)
					return
				}
			}

			if tt.resp != nil {
				var r ResponseSingle
				err = json.Unmarshal(body, &r)
				if err != nil {
					t.Errorf("invalid response body: %v", err)
					return
				}

				a := testArticle(r.Article) // Coerse to custom type with Equals method
				if !a.Equals(&tt.resp.Article) {
					t.Errorf("response not matched, expected '%+v', got '%+v'", tt.resp.Article, r.Article)
				}
			}
		})
	}
}
