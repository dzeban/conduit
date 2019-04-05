package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/dzeban/conduit/app"
)

var server *Server

func TestMain(m *testing.M) {
	mockConf := Config{
		Secret: "mock",
		Articles: app.ArticlesConfig{
			Type: "mock",
		},
		Users: app.UsersConfig{
			Type: "mock",
		},
	}

	var err error
	server, err = NewServer(mockConf)
	if err != nil {
		panic("failed to create server: " + err.Error())
	}

	os.Exit(m.Run())
}

func TestArticles(t *testing.T) {
	tests := []struct {
		name   string
		method string
		url    string
		status int
		body   string
	}{
		{"list", "GET", "/articles/", 200, "Title 1"},
		{"notfound", "GET", "/articles/xxx", 404, ""},
		{"single", "GET", "/articles/slug-2", 200, "Description 2"},
	}

	for _, expected := range tests {
		// Run in a subtest to distinguish tests by name
		t.Run(expected.name, func(t *testing.T) {
			req := httptest.NewRequest(expected.method, expected.url, nil)

			rr := httptest.NewRecorder()
			server.httpServer.Handler.ServeHTTP(rr, req)

			// Check status
			status := rr.Code
			if status != expected.status {
				t.Errorf("invalid status code: expected %v got %v'", expected.status, status)
			}

			// Check body
			body := rr.Body.String()
			if !strings.Contains(body, expected.body) {
				t.Errorf("invalid body: expected %v, got %v", expected.body, body)
			}
		})
	}
}

func TestUserRegister(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		reqBody  string
		status   int
		response interface{}
	}{
		{
			"Register",
			"POST",
			"/users",
			`{"user":{"username":"user3","email":"user3@example.com","password":"user3pass"}}`,
			201,
			app.UserRequest{
				User: app.User{
					Name:  "user3",
					Email: "user3@example.com",
					Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWV9.LCCSzgQvBNx6xE8P2xJurQ_ykszQIDqyRDL28AeBCls`,
				},
			},
		},
		{
			"RegisterExisting",
			"POST",
			"/users",
			`{"user":{"username":"user3","email":"user3@example.com","password":"user3pass"}}`,
			422,
			nil,
		},
		{
			"InvalidNoEmail",
			"POST",
			"/users",
			`{"user":{"username":"noemailuser","password":"noemailuserpass"}}`,
			422,
			nil,
		},
		{
			"InvalidNoUsername",
			"POST",
			"/users",
			`{"user":{"email":"nousername@example.com","password":"nousernamepass"}}`,
			422,
			nil,
		},
		{
			"InvalidNoPassword",
			"POST",
			"/users",
			`{"user":{"email":"nopassworduser@example.com","username":"nopassworduser"}}`,
			422,
			nil,
		},
	}

	for _, expected := range tests {
		// Run in a subtest to distinguish tests by name
		t.Run(expected.name, func(t *testing.T) {
			req := httptest.NewRequest(expected.method, expected.url, strings.NewReader(expected.reqBody))

			rr := httptest.NewRecorder()
			server.httpServer.Handler.ServeHTTP(rr, req)

			// Check status
			status := rr.Code
			if status != expected.status {
				t.Errorf("invalid status code: expected %v got %v", expected.status, status)
			}

			// fmt.Println(rr.Body.String())

			// Check non-error response
			if expected.status >= 200 && expected.status < 400 {
				var user app.UserRequest
				err := json.Unmarshal(rr.Body.Bytes(), &user)
				if err != nil {
					t.Error("failed to unmarshal json", err)
				}

				if !reflect.DeepEqual(expected.response, user) {
					t.Errorf("users not matching: expected %v got %v", expected.response, user)
				}
			}
		})
	}
}

func TestUserRegisterToken(t *testing.T) {
	userRegisterRequest := `{"user": {"username":"aaa","email":"a@example.com","password":"123"}}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(userRegisterRequest))
	rr := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusCreated {
		t.Errorf("invalid status code: expected %v got %v", http.StatusCreated, status)
	}

	type RespToken struct {
		User struct {
			Token string `json:"token"`
		} `json:"user"`
	}
	var respToken RespToken

	resp := rr.Body.Bytes()
	err := json.Unmarshal(resp, &respToken)
	if err != nil {
		t.Error("failed to unmarshal JSON response", err)
	}

	expectedJWT := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWV9.LCCSzgQvBNx6xE8P2xJurQ_ykszQIDqyRDL28AeBCls`
	if respToken.User.Token != expectedJWT {
		t.Errorf("jwt not expected: expected %v got %v", expectedJWT, respToken.User.Token)
		t.Errorf("req was %v, resp was %v", userRegisterRequest, rr.Body.String())
	}
}

func TestUserRegisterNoPassword(t *testing.T) {
	userRegisterRequest := `{"user": {"username":"boss","email":"boss@example.com","password":"bossypassword"}}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(userRegisterRequest))
	rr := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusCreated {
		t.Errorf("invalid status code: expected %v got %v", http.StatusCreated, status)
	}

	resp := rr.Body.String()
	if strings.Contains(resp, "bossypassword") {
		t.Errorf("plaintext password found in response %#v", rr.Body.String())
	}
}
