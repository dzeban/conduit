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
			http.StatusCreated,
			app.UserRequest{
				User: app.User{
					Name:  "user3",
					Email: "user3@example.com",
					Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InVzZXIzQGV4YW1wbGUuY29tIn0.e3e-9S3ejZpyUgLcfDWuTIGUR79P5I6-xk_cu0RXks0`,
				},
			},
		},
		{
			"RegisterExisting",
			"POST",
			"/users",
			`{"user":{"username":"user3","email":"user3@example.com","password":"user3pass"}}`,
			http.StatusConflict,
			nil,
		},
		{
			"InvalidNoEmail",
			"POST",
			"/users",
			`{"user":{"username":"noemailuser","password":"noemailuserpass"}}`,
			http.StatusBadRequest,
			nil,
		},
		{
			"InvalidNoUsername",
			"POST",
			"/users",
			`{"user":{"email":"nousername@example.com","password":"nousernamepass"}}`,
			http.StatusBadRequest,
			nil,
		},
		{
			"InvalidNoPassword",
			"POST",
			"/users",
			`{"user":{"email":"nopassworduser@example.com","username":"nopassworduser"}}`,
			http.StatusBadRequest,
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

	expectedJWT := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6ImFAZXhhbXBsZS5jb20ifQ.fgpEq8igZEIW7tePDA6CAFk7OY8zt1q54-Sxx-kZNBg`
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

func TestUserLogin(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		url      string
		reqBody  string
		status   int
		response interface{}
	}{
		{
			"Login",
			"POST",
			"/users/login",
			`{"user":{"email":"user1@example.com","password":"user1pass"}}`,
			http.StatusOK,
			app.UserRequest{
				User: app.User{
					Name:  "user1",
					Email: "user1@example.com",
					Token: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InVzZXIxQGV4YW1wbGUuY29tIn0.hZzfk_8yUXuXPu9B4lB4sbh04L4rfxC9Rmqf22HMGX8`,
				},
			},
		},
		{
			"LoginUnauthorized",
			"POST",
			"/users/login",
			`{"user":{"email":"nosuchuser@example.com","password":"nosuchuserpassword"}}`,
			http.StatusUnauthorized,
			nil,
		},
		{
			"InvalidNoEmail",
			"POST",
			"/users/login",
			`{"user":{"password":"noemailuserpass"}}`,
			http.StatusBadRequest,
			nil,
		},
		{
			"InvalidNoPassword",
			"POST",
			"/users",
			`{"user":{"email":"nopassworduser@example.com"}}`,
			http.StatusBadRequest,
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

func TestJWTParse(t *testing.T) {
	tests := []struct {
		name   string
		header string
		secret []byte
		isErr  bool
	}{
		{
			"Successful",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.8dbE6VKua4RQsqJZXTqisRXCBf4K5dTgBgPYHwn1ikc",
			[]byte(""),
			false,
		},
		{
			"InvalidFormat",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.8dbE6VKua4RQsqJZXTqisRXCBf4K5dTgBgPYHwn1ikc",
			[]byte(""),
			true,
		},
		{
			"InvalidFormatName",
			"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.8dbE6VKua4RQsqJZXTqisRXCBf4K5dTgBgPYHwn1ikc",
			[]byte(""),
			true,
		},
		{
			"InvalidJWT",
			"Token ZZZZZZZZZZZZZZZZNiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.8dbE6VKua4RQsqJZXTqisRXCBf4K5dTgBgPYHwn1ikc",
			[]byte(""),
			true,
		},
		{
			"InvalidSignature",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InRlc3RAZXhhbXBsZS5jb20ifQ.8dbE6VKua4RQsqJZXTqisRXCBf4K5dTgBgPYHwn1ik1",
			[]byte(""),
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseJWTClaimsFromHeader(test.header, test.secret)

			if (err != nil) != test.isErr {
				t.Errorf("error is expected to be %v, got %v, header is '%v', secret is '%s'\n", err != nil, test.isErr, test.header, test.secret)
			}
		})
	}
}

func TestJWTAuth(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		status     int
		body       string
	}{
		{
			"NoHeader",
			"",
			http.StatusUnauthorized,
			ErrJWTNoAuthorizationHeader.Error(),
		},
		{
			"InvalidHeader",
			"bebebe",
			http.StatusBadRequest,
			"",
		},
		{
			"NoSignedClaim",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0QGV4YW1wbGUuY29tIn0.TLTUpuQq8FH-JaUfnho9dkSB_XKTlDCxAdiLsMJ-TdA",
			http.StatusUnauthorized,
			ErrJWTNoSignedClaim.Error(),
		},
		{
			"NoSubClaim",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWV9.LCCSzgQvBNx6xE8P2xJurQ_ykszQIDqyRDL28AeBCls",
			http.StatusUnauthorized,
			ErrJWTNoSubClaim.Error(),
		},
	}

	for _, expected := range tests {
		t.Run(expected.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users", nil)
			if expected.authHeader != "" {
				req.Header["Authorization"] = []string{expected.authHeader}
			}

			rr := httptest.NewRecorder()

			server.httpServer.Handler.ServeHTTP(rr, req)

			// Check status
			status := rr.Code
			if status != expected.status {
				t.Errorf("invalid status code: expected %v got %v", expected.status, status)
			}

			if expected.body != "" {
				body := rr.Body.String()
				if !strings.Contains(body, expected.body) {
					t.Errorf("unexpected error message: expected '%v', got '%v'\n", expected.body, body)
				}
			}
		})
	}
}

func TestUserGet(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		status     int
		body       string
	}{
		{
			"Successful",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InVzZXIxQGV4YW1wbGUuY29tIn0.hZzfk_8yUXuXPu9B4lB4sbh04L4rfxC9Rmqf22HMGX8",
			http.StatusOK,
			"",
		},
		{
			"InvalidSub",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6MTIzfQ.kjCfjIKA_buC-tNby6aeh5cklId7J1qWj0qn6rcDAP0",
			http.StatusUnauthorized,
			"invalid auth email",
		},
		{
			"NoUser",
			"Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWduZWQiOnRydWUsInN1YiI6InBpcGthQGV4YW1wbGUuY29tIn0.i--wUiS2g7XPrL83EPo5E_8S3vGh58RRl3AKAZnz8j0",
			http.StatusNotFound,
			app.ErrUserNotFound.Error(),
		},
	}

	for _, expected := range tests {
		t.Run(expected.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users", nil)
			if expected.authHeader != "" {
				req.Header["Authorization"] = []string{expected.authHeader}
			}

			rr := httptest.NewRecorder()

			server.httpServer.Handler.ServeHTTP(rr, req)

			// Check status
			status := rr.Code
			if status != expected.status {
				t.Errorf("invalid status code: expected %v got %v", expected.status, status)
			}

			if expected.body != "" {
				body := rr.Body.String()
				if !strings.Contains(body, expected.body) {
					t.Errorf("unexpected error message: expected '%v', got '%v'\n", expected.body, body)
				}
			}
		})
	}
}
