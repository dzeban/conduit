package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
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
			http.StatusUnprocessableEntity,
		},
		{
			"empty",
			"{}",
			http.StatusUnprocessableEntity,
		},
		{
			"invalid",
			`{"user": {"email":"test@example.com"}}`,
			http.StatusUnprocessableEntity,
		},
		{
			"incorrect",
			`{"user":{"email":"test@example.com","password":"incorrect"}}`,
			http.StatusUnprocessableEntity,
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

	url := ts.URL + "/login"

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
			http.StatusUnprocessableEntity,
		},
		{
			"empty",
			"{}",
			http.StatusUnprocessableEntity,
		},
		{
			"conflict",
			`{"user":{"email":"test@example.com","username": "test", "password":"test"}}`,
			http.StatusUnprocessableEntity,
		},
		{
			"invalid",
			`{"user":{"email":"new@example.com","password":"new"}}`,
			http.StatusUnprocessableEntity,
		},
		{
			"valid",
			`{"user":{"email":"new@example.com","username": "new", "password":"new"}}`,
			http.StatusCreated,
		},
	}

	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", ts.URL, strings.NewReader(c.body))
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

	_, err = jwt.Parse(r.User.Token, []byte(testSecret))
	if err != nil {
		t.Errorf("invalid token '%s': %s", r.User.Token, err)
	}
}

func TestHandleUserGetEmpty(t *testing.T) {
	cases := []struct {
		name   string
		auth   string
		status int
	}{
		{
			"null",
			"",
			http.StatusUnauthorized,
		},
		{
			"empty",
			"Token ",
			http.StatusUnauthorized,
		},
	}
	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL, nil)
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

func TestHandleUserGetJWT(t *testing.T) {
	cases := []struct {
		name   string
		user   app.User
		status int
	}{
		{
			"noemail",
			app.User{Name: "test"},
			http.StatusUnauthorized,
		},
		{
			"noname",
			app.User{Email: "test@example.com"},
			http.StatusUnauthorized,
		},
		{
			"emptyemail",
			app.User{Email: ""},
			http.StatusUnauthorized,
		},
		{
			"notfound",
			app.User{Name: "nosuchuser", Email: "nosuchuser@example.com"},
			http.StatusNotFound,
		},
		{
			"valid",
			app.User{Name: "test", Email: "test@example.com"},
			http.StatusOK,
		},
	}

	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL, nil)
			if err != nil {
				t.Errorf("failed to create request: %s", err)
			}
			token, err := jwt.New(&c.user, []byte(testSecret))
			if err != nil {
				t.Errorf("failed to make jwt: %v", err)
			}

			req.Header.Add("Authorization", "Token "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("failed to make a request: %s", err)
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

func TestHandleUserUpdate(t *testing.T) {
	cases := []struct {
		name   string
		user   app.User
		body   string
		status int
	}{
		{
			"valid",
			app.User{Name: "test", Email: "test@example.com"},
			`{"user": {"username": "admin"}}`,
			http.StatusOK,
		},
	}

	s := New(newMockStore(), testSecret)

	ts := httptest.NewServer(s.router)
	defer ts.Close()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", ts.URL, strings.NewReader(c.body))
			if err != nil {
				t.Fatalf("failed to create request: %s", err)
			}

			token, err := jwt.New(&c.user, []byte(testSecret))
			if err != nil {
				t.Errorf("failed to make jwt: %v", err)
			}

			req.Header.Add("Authorization", "Token "+token)

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

			if resp.StatusCode == http.StatusOK {
				checkToken(body, t)
			}
		})
	}
}
