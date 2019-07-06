// +build integration

package user

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dzeban/conduit/app"
)

var (
	service *Service
	server  *httptest.Server
	user    *app.User
)

const (
	DSN    = "postgres://test:test@postgres:5432/test?sslmode=disable"
	secret = "test"
)

// TestMain initializes integration test environment by creating service
// instance with a predefined user account. It also create test server used for
// testing http handlers.
func TestMain(m *testing.M) {
	var err error
	service, err = NewService(DSN, secret)
	if err != nil {
		panic("failed to create service: " + err.Error())
	}

	server = httptest.NewServer(service.router)
	defer server.Close()

	// Prepare test user
	testUser := app.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: "test",
	}

	user, err = service.Register(app.UserRequest{User: testUser})
	if err != nil {
		panic("failed to prepare test user: " + err.Error())
	}

	os.Exit(m.Run())
}
