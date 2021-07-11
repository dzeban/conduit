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

// httpTest describe tests in the current file
type httpTest struct {
	name   string          // Test name
	user   *app.User       // App user stored in the context that is making request
	req    *http.Request   // Test request to send
	status int             // Expected HTTP request status code
	resp   *ResponseSingle // Expected response unmarshalled from JSON
	err    error           // Expected error
}

// check fails test when httpTest tt is not passed. It validates HTTP response
// status code, check for expected error and compare returned response with
// expected.
func check(t *testing.T, tt httpTest, resp *http.Response) {
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
		err := json.Unmarshal(body, &r)
		if err != nil {
			t.Errorf("invalid response body: %v", err)
			return
		}

		a := testArticle(r.Article) // Coerse to custom type with Equals method
		if !a.Equals(&tt.resp.Article) {
			t.Errorf("response not matched, expected '%+v', got '%+v'", tt.resp.Article, r.Article)
		}
	}
}

func TestHTTPHandlers(t *testing.T) {
	tests := []httpTest{
		// --- Create handler ---
		{
			"Create/Null",
			&app.User{},
			httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")),
			http.StatusUnauthorized,
			nil,
			// Don't check concrete error because it is returned by JWT middleware.
			// We need to be sure that unauthorized request returns error.
			nil,
		},
		{
			"Create/Unmarshal",
			&mock.UserValid,
			httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{")),
			http.StatusUnprocessableEntity,
			nil,
			errorInvalidRequest,
		},
		{
			"Create/InvalidRequest",
			&mock.UserValid,
			httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"article":{"title":"","description":"some","body":"some"}}`),
			),
			http.StatusUnprocessableEntity,
			nil,
			errorValidationTitleIsRequired,
		},
		{
			"Create/Valid",
			&mock.UserValid,
			httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(`{"article":{"title":"new","description":"new","body":"new"}}`),
			),
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

		// --- Delete handler ---
		{
			"Delete/NoUser",
			&app.User{},
			httptest.NewRequest(http.MethodDelete, "/"+mock.ArticleValid.Slug, nil),
			http.StatusUnauthorized,
			nil,
			// Don't check concrete error because it is returned by JWT middleware.
			// We need to be sure that unauthorized request returns error.
			nil,
		},
		{
			"Delete/NonExisting",
			&mock.UserValid,
			httptest.NewRequest(http.MethodDelete, "/xxxxxx", nil),
			http.StatusUnprocessableEntity,
			nil,
			errorArticleNotFound,
		},
		{
			"Delete/Valid",
			&mock.UserValid,
			httptest.NewRequest(http.MethodDelete, "/"+mock.ArticleToDelete.Slug, nil),
			http.StatusOK,
			nil,
			nil,
		},

		// --- Update handler ---
		{
			"Update/NoUser",
			&app.User{},
			httptest.NewRequest(http.MethodPut, "/"+mock.ArticleValid.Slug, strings.NewReader("{}")),
			http.StatusUnauthorized,
			nil,
			// Don't check concrete error because it is returned by JWT middleware.
			// We need to be sure that unauthorized request returns error.
			nil,
		},
		{
			"Update/Unmarshal",
			&mock.UserValid,
			httptest.NewRequest(http.MethodPut, "/"+mock.ArticleValid.Slug, strings.NewReader("{")),
			http.StatusUnprocessableEntity,
			nil,
			errorInvalidRequest,
		},
		{
			"Update/NonExisting",
			&mock.UserValid,
			httptest.NewRequest(http.MethodPut, "/xxxxx", strings.NewReader(`{"article": {"body":"test new body"}}`)),
			http.StatusUnprocessableEntity,
			nil,
			errorArticleNotFound,
		},
		{
			"Update/Valid",
			&mock.UserValid,
			httptest.NewRequest(http.MethodPut, "/"+mock.ArticleUpdated.Slug, strings.NewReader(`{"article":{"title":"test new title"}}`)),
			http.StatusOK,
			&ResponseSingle{
				app.Article{
					Title:       "test new title",
					Slug:        "other-title-azxs",
					Description: "Other description",
					Body:        "Other body",
					Author:      mock.Profile1,
				},
			},
			nil,
		},

		// --- Get handler ---
		{
			"Get/NonExisting",
			nil,
			httptest.NewRequest(http.MethodGet, "/xxxxx", nil),
			http.StatusUnprocessableEntity,
			&ResponseSingle{app.Article{}}, // Ensure that nothing was returned in case of not found
			errorArticleNotFound,
		},
		{
			"Get/Valid",
			nil,
			httptest.NewRequest(http.MethodGet, "/"+mock.ArticleValid.Slug, nil),
			http.StatusOK,
			&ResponseSingle{mock.ArticleValid},
			nil,
		},
	}

	s, err := NewHTTP(mock.NewArticleStore(), mock.NewProfilesStore(), []byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwt.New(tt.user, []byte(testSecret))
			if err != nil {
				t.Errorf("failed to make JWT")
				return
			}

			tt.req.Header.Add("Authorization", "Token "+token)
			rr := httptest.NewRecorder()

			s.ServeHTTP(rr, tt.req)
			resp := rr.Result()
			check(t, tt, resp)
		})
	}
}

// func TestListHandler(t *testing.T) {
// 	tests := []httpTest{
// 		// LimitValidation
// 		// OffsetValidation
// 		// ByAuthor
// 		// Public
// 		{
// 			"NonExisting",
// 			nil,
// 			"/xxxxxx",
// 			"",
// 			http.StatusUnprocessableEntity,
// 			&ResponseSingle{app.Article{}}, // Ensure that nothing was returned in case of not found
// 			errorArticleNotFound,
// 		},
// 		{
// 			"Valid",
// 			nil,
// 			"/" + mock.ArticleValid.Slug,
// 			"",
// 			http.StatusOK,
// 			&ResponseSingle{mock.ArticleValid},
// 			nil,
// 		},
// 	}

// 	s, err := NewHTTP(mock.NewArticleStore(), mock.NewProfilesStore(), []byte(testSecret))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			token, err := jwt.New(tt.user, []byte(testSecret))
// 			if err != nil {
// 				t.Errorf("failed to make JWT")
// 				return
// 			}

// 			req := httptest.NewRequest(http.MethodGet, tt.target, strings.NewReader(tt.body))
// 			req.Header.Add("Authorization", "Token "+token)

// 			rr := httptest.NewRecorder()

// 			s.ServeHTTP(rr, req)

// 			resp := rr.Result()

// 			check(t, tt, resp)
// 		})
// 	}
// }
