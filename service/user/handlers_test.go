// +build integration

package user

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

//
// XXX: Integration test environment initialized in service_test.go
//

func TestHandleUserLogin(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		status int
		// TODO: body string - compare serialized response
	}{
		{
			name:   "null",
			data:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "empty",
			data:   "{}",
			status: http.StatusBadRequest,
		},
		{
			name:   "invalid",
			data:   `{"user": {"email":"test@example.com"}}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "incorrect",
			data:   `{"user":{"email":"test@example.com","password":"incorrect"}}`,
			status: http.StatusUnauthorized,
		},
		{
			name:   "valid",
			data:   `{"user":{"email":"test@example.com","password":"test"}}`,
			status: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := http.Post(server.URL+"/users/login", "application/json", strings.NewReader(test.data))
			if err != nil {
				t.Error("failed to make a request:", err)
			}
			defer resp.Body.Close()

			if test.status != resp.StatusCode {
				t.Errorf("invalid status code: expected %d, got %d", test.status, resp.StatusCode)
				body, _ := ioutil.ReadAll(resp.Body)
				t.Error(string(body))
			}
		})
	}
}

func TestHandleUserRegister(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		status int
		// TODO: body string - compare serialized response
	}{
		{
			name:   "null",
			data:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "empty",
			data:   "{}",
			status: http.StatusBadRequest,
		},
		{
			name:   "conflict",
			data:   `{"user":{"email":"test@example.com","username": "test", "password":"test"}}`,
			status: http.StatusConflict,
		},
		{
			name:   "invalid",
			data:   `{"user":{"email":"newuser@example.com","password":"newuser"}}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "valid",
			data:   `{"user":{"email":"newuser@example.com","username": "newuser", "password":"newuser"}}`,
			status: http.StatusCreated,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := http.Post(server.URL+"/users/", "application/json", strings.NewReader(test.data))
			if err != nil {
				t.Error("failed to make a request:", err)
			}
			defer resp.Body.Close()

			if test.status != resp.StatusCode {
				t.Errorf("invalid status code: expected %d, got %d", test.status, resp.StatusCode)
				body, _ := ioutil.ReadAll(resp.Body)
				t.Error(string(body))
			}
		})
	}
}
