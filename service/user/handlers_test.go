// +build integration

package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

//
// XXX: Integration test environment initialized in service_test.go
//

// userTest is a type to check response from handlers.
// We don't use app.User and app.UserRequest type to catch the case when app.User is changed but all
// other clients are expecting the previous version.
type userTest struct {
	Name  string `json:"username"`
	Email string `json:"email"`
	Bio   string `json:"bio,omitempty"`
	Image string `json:"image,omitempty"` // base64 encoded
	Token string `json:"token,omitempty"`
}

type userResponse struct {
	User userTest `json:"user"`
}

func TestHandleUserLogin(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		status   int
		response *userResponse
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
			response: &userResponse{
				User: userTest{
					Name:  "test",
					Email: "test@example.com",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := http.Post(server.URL+"/users/login", "application/json", strings.NewReader(test.data))
			if err != nil {
				t.Error("failed to make a request:", err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			if test.status != resp.StatusCode {
				t.Errorf("invalid status code: expected %d, got %d", test.status, resp.StatusCode)
				t.Error(string(body))
			}

			if test.response != nil {
				err = checkResponse(test.response, body)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestHandleUserRegister(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		status   int
		response *userResponse
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
			response: &userResponse{
				User: userTest{
					Name:  "newuser",
					Email: "newuser@example.com",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := http.Post(server.URL+"/users/", "application/json", strings.NewReader(test.data))
			if err != nil {
				t.Error("failed to make a request:", err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			if test.status != resp.StatusCode {
				t.Errorf("invalid status code: expected %d, got %d", test.status, resp.StatusCode)
			}

			if test.response != nil {
				err = checkResponse(test.response, body)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func checkResponse(expected *userResponse, body []byte) error {
	// Check response
	var response userResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %s", err)
	}

	// Check simple fields
	if expected.User.Name != response.User.Name {
		return fmt.Errorf("invalid user: expected %v, got %v", expected.User.Name, response.User.Name)
	}

	if expected.User.Email != response.User.Email {
		return fmt.Errorf("invalid email: expected %v, got %v", expected.User.Email, response.User.Email)
	}

	// Check token by parsing it
	token, err := jwt.Parse(response.User.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if !token.Valid {
		return fmt.Errorf("invalid token: %s", err)
	}

	return nil
}
