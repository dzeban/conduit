package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

const testSecret = "test"

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		name   string
		body   string
		status int
	}{
		{
			"null",
			"",
			http.StatusBadRequest,
		},
		{
			"empty",
			"{}",
			http.StatusBadRequest,
		},
		{
			"invalid",
			`{"user": {"email":"test@example.com"}}`,
			http.StatusBadRequest,
		},
		{
			"incorrect",
			`{"user":{"email":"test@example.com","password":"incorrect"}}`,
			http.StatusUnauthorized,
		},
		{
			"valid",
			`{"user":{"email":"test@example.com","password":"test"}}`,
			http.StatusOK,
		},
	}

	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	url := ts.URL + "/users/login"

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", url, strings.NewReader(c.body))
			if err != nil {
				t.Fatalf("failed to create request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to make a request: %s", err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != c.status {
				t.Errorf("unexpected status, want %d, got %d", c.status, resp.StatusCode)
				t.Error(string(body))
			}

			if resp.StatusCode == http.StatusOK {
				checkToken(body, t)
			}
		})
	}
}

func TestHandleUserRegister(t *testing.T) {
	cases := []struct {
		name   string
		body   string
		status int
	}{
		{
			"null",
			"",
			http.StatusBadRequest,
		},
		{
			"empty",
			"{}",
			http.StatusBadRequest,
		},
		{
			"conflict",
			`{"user":{"email":"test@example.com","username": "test", "password":"test"}}`,
			http.StatusConflict,
		},
		{
			"invalid",
			`{"user":{"email":"new@example.com","password":"new"}}`,
			http.StatusBadRequest,
		},
		{
			"valid",
			`{"user":{"email":"new@example.com","username": "new", "password":"new"}}`,
			http.StatusCreated,
		},
	}

	s := New(newMockStore(), testSecret)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/", strings.NewReader(c.body))
			s.HandleUserRegister(rr, req)

			if rr.Code != c.status {
				t.Errorf("unexpected status, want %d, got %d", c.status, rr.Code)
				t.Error(rr.Body.String())
			}

			if rr.Code == http.StatusOK {
				checkToken(rr.Body.Bytes(), t)
			}
		})
	}
}

func checkToken(body []byte, t *testing.T) {
	type resp struct {
		User struct {
			Token string `json:"token"`
		} `json:"user"`
	}

	var r resp
	err := json.Unmarshal(body, &r)
	if err != nil {
		t.Errorf("failed to unmarshal json: %s", err)
	}

	_, err = jwt.Parse(r.User.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})

	if err != nil {
		t.Errorf("invalid token '%s': %s", r.User.Token, err)
	}
}

func TestHandleUserGet(t *testing.T) {
	cases := []struct {
		name   string
		auth   string
		status int
	}{
		{
			"null",
			"",
			http.StatusBadRequest,
		},
		{
			"empty",
			"Token ",
			http.StatusBadRequest,
		},
		{
			// no sub claim
			"nosub",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWV9.6ARuTLidiCvLg5nLJhrWff9fLbZaQTvRKKBQW-04P9Y",
			http.StatusUnauthorized,
		},
		{
			// empty sub claim
			"emptysub",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6IiJ9.R7UsDbYl0wVvAate0SbP8nDdXBp3uOVF-gP8FaegaZg",
			http.StatusUnauthorized,
		},
		{
			// email is nosuchuser@example.com
			"notfound",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6Im5vc3VjaHVzZXJAZXhhbXBsZS5jb20ifQ.7Ckyqr4bsJRSSsEjRcNmskSNqhhPQkqBi2huaFX9MRY",
			http.StatusNotFound,
		},
		{
			// email is test@example.com
			"valid",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.eox2GVi27V5h16i_Ob5KtEnOtiMBu-jzpapDdeYzFbI",
			http.StatusOK,
		},
	}

	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL+"/users/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %s", err)
			}
			req.Header.Add("Authorization", c.auth)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to make a request: %s", err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != c.status {
				t.Errorf("invalid status code: want %d, got %d", c.status, resp.StatusCode)
				t.Error(string(body))
			}
		})
	}
}

// // registerTestUser registers new user and obtains token
// func registerTestUser(email, username, password string) (*userResponse, error) {
// 	userData := fmt.Sprintf(`{"user":{"email":"%s","username": "%s", "password":"%s"}}`, email, username, password)
// 	resp, err := http.Post(testServer.URL+"/users/", "application/json", strings.NewReader(userData))
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to make user register request")
// 	}

// 	if resp.StatusCode != http.StatusCreated {
// 		return nil, fmt.Errorf("invalid status code: expected %d, got %d", http.StatusCreated, resp.StatusCode)
// 	}

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()

// 	var response userResponse
// 	err = json.Unmarshal(body, &response)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to unmarshal response")
// 	}

// 	return &response, nil
// }
