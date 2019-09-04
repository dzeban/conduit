package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsIm5hbWUiOiJ0ZXN0In0.nZMAfsqMBNKSc7zD_F45icTTMolVMARBGOK13INJdtw",
			http.StatusUnauthorized,
		},
		{
			// no name claim
			"noname",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.eox2GVi27V5h16i_Ob5KtEnOtiMBu-jzpapDdeYzFbI",
			http.StatusUnauthorized,
		},
		{
			// empty sub claim
			"emptysub",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6IiIsIm5hbWUiOiJ0ZXN0In0.24fNqWxFncXeBVi4gMk6wQJ9iSMrxZ9_CgvFNG8djno",
			http.StatusUnauthorized,
		},
		{
			// email is nosuchuser@example.com
			"notfound",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6Im5vc3VjaHVzZXJAZXhhbXBsZS5jb20iLCJuYW1lIjoibm9zdWNodXNlciJ9.fPIrYSf8RF8rp_oI5RjkY68ex-mIz87erD0SqCiHR7I",
			http.StatusNotFound,
		},
		{
			// email is test@example.com
			"valid",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20iLCJuYW1lIjoidGVzdCJ9.2wA__EfTfnZ2LEDUz3cB1lxpYWWo9w3THE-fa0LqzaU",
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

func TestHandleUserUpdate(t *testing.T) {
	cases := []struct {
		name   string
		auth   string
		body   string
		status int
	}{
		{
			"null",
			"",
			`{}`,
			http.StatusBadRequest,
		},
		{
			"empty",
			"Token ",
			`{}`,
			http.StatusBadRequest,
		},
		{
			// no sub claim
			"nosub",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWV9.6ARuTLidiCvLg5nLJhrWff9fLbZaQTvRKKBQW-04P9Y",
			`{}`,
			http.StatusUnauthorized,
		},
		{
			// empty sub claim
			"emptysub",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6IiJ9.R7UsDbYl0wVvAate0SbP8nDdXBp3uOVF-gP8FaegaZg",
			`{}`,
			http.StatusUnauthorized,
		},
		{
			"invalid",
			// email is test@example.com
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20iLCJuYW1lIjoidGVzdCJ9.2wA__EfTfnZ2LEDUz3cB1lxpYWWo9w3THE-fa0LqzaU",
			`{}`,
			http.StatusBadRequest,
		},
		{
			"other",
			// email is test@example.com
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20iLCJuYW1lIjoidGVzdCJ9.2wA__EfTfnZ2LEDUz3cB1lxpYWWo9w3THE-fa0LqzaU",
			`{"user": {"email": "other@example.com", "password":"evil"}}`,
			http.StatusForbidden,
		},
		{
			"valid",
			// email is test@example.com
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20iLCJuYW1lIjoidGVzdCJ9.2wA__EfTfnZ2LEDUz3cB1lxpYWWo9w3THE-fa0LqzaU",
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

			if resp.StatusCode == http.StatusOK {
				checkToken(body, t)
			}
		})
	}
}
