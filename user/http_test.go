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
	"github.com/dzeban/conduit/mock"
	"github.com/dzeban/conduit/transport"
)

const testSecret = "test"

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
		err    error
	}{
		{
			"Null",
			"",
			http.StatusUnprocessableEntity,
			transport.ErrorUnmarshal,
		},
		{
			"Empty",
			"{}",
			http.StatusUnauthorized,
			nil, // Don't check for specific validation error because validation order may change
		},
		{
			"IncorrectPassword",
			`{"user":{"email":"test@example.com","password":"incorrect"}}`,
			http.StatusUnauthorized,
			app.ErrorLogin,
		},
		{
			"valid",
			`{"user":{"email":"test@example.com","password":"test"}}`,
			http.StatusOK,
			nil,
		},
	}

	s, err := NewHTTP(mock.NewUserStore(), []byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
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
		})
	}
}

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name        string
		contextUser *app.User
		status      int
		err         error
	}{
		{
			"Unauthorized",
			&app.User{},
			http.StatusUnauthorized,
			// Don't check concrete error because it is returned by JWT middleware.
			// We need to be sure that unauthorized request returns error.
			nil,
		},
		{
			"NotFound",
			&app.User{
				Name:  "no_such_user",
				Email: "no_such_user@example.com",
			},
			http.StatusUnprocessableEntity,
			app.ErrorUserNotFound,
		},
		{
			"Valid",
			&app.User{
				Name:  "test",
				Email: "test@example.com",
			},
			http.StatusOK,
			nil,
		},
	}

	s, err := NewHTTP(mock.NewUserStore(), []byte(testSecret))
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

			req := httptest.NewRequest(http.MethodGet, "/", nil)
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
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
		err    error
	}{
		{
			"Null",
			"",
			http.StatusUnprocessableEntity,
			transport.ErrorUnmarshal,
		},
		{
			"Empty",
			"{}",
			http.StatusUnprocessableEntity,
			nil, // Don't check for specific validation error because validation order may change
		},
		{
			"NoUsername",
			`{"user":{"email":"test@example.com","password":"test"}}`,
			http.StatusUnprocessableEntity,
			app.ErrorRegister,
		},
		{
			"valid",
			`{"user":{"username": "new_register", "email":"new_register@example.com","password":"new_register"}}`,
			http.StatusCreated,
			nil,
		},
	}

	s, err := NewHTTP(mock.NewUserStore(), []byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
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
		})
	}

}

func TestLoginRegisterHandlerToken(t *testing.T) {
	// Check that JWT obtained during register and login is valid by doing
	// getting current user.
	tests := []struct {
		name      string
		req       string
		path      string
		expStatus int
		exp       string
	}{
		{
			"Login",
			`{"user":{"email":"test@example.com","password":"test"}}`,
			"/login",
			http.StatusOK,
			"test",
		},
		{
			"Register",
			`{"user":{"username": "new_register", "email":"new_register@example.com","password":"new_register"}}`,
			"/",
			http.StatusCreated,
			"new_register",
		},
	}

	s, err := NewHTTP(mock.NewUserStore(), []byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		// Perform request to Login or Register handler
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.req))
		s.ServeHTTP(rr, req)

		res := rr.Result()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != tt.expStatus {
			t.Fatalf("incorrect status during %v user: expected %v, got %v, body is %v", tt.name, tt.expStatus, res.StatusCode, string(body))
		}

		var resp Response
		if err = json.Unmarshal(body, &resp); err != nil {
			t.Fatal(err)
		}

		// Grab JWT
		token := resp.User.Token

		// Try to get current user
		rr = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "Token "+token)

		s.ServeHTTP(rr, req)

		res = rr.Result()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("incorrect status during Get user: expected %v, got %v", http.StatusOK, res.StatusCode)
		}

		if !strings.Contains(string(body), tt.exp) {
			t.Errorf(`expected string not found in Get user: expected %v', body is '%v', token is '%v'`, tt.exp, string(body), token)
		}
	}
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
		err    error
	}{
		{
			"Null",
			"",
			http.StatusUnprocessableEntity,
			transport.ErrorUnmarshal,
		},
		{
			"Empty",
			"{}",
			http.StatusUnprocessableEntity,
			nil, // Don't check for specific validation error because validation order may change
		},
		{
			"Forbidden",
			`{"user":{"email":"updated@example.com","password":"test"}}`,
			http.StatusUnauthorized,
			app.ErrorUserUpdateForbidden,
		},
		{
			"Valid",
			`{"user":{"email":"test@example.com","password":"test"}}`,
			http.StatusOK,
			nil,
		},
	}

	s, err := NewHTTP(mock.NewUserStore(), []byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	token, err := jwt.New(&mock.UserValid, []byte(testSecret))
	if err != nil {
		t.Fatal("failed to make JWT")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(tt.body))
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
		})
	}

}
