// +build integration

package user

import (
	"os"
	"testing"

	"github.com/go-test/deep"

	"github.com/dzeban/conduit/app"
)

var (
	service *Service
	user    *app.User
)

const (
	DSN    = "postgres://test:test@postgres:5432/test?sslmode=disable"
	secret = "test"
)

func TestMain(m *testing.M) {
	var err error
	service, err = NewService(DSN, secret)
	if err != nil {
		panic("failed to create service: " + err.Error())
	}

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

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		email string
		user  *app.User
		err   error
	}{
		{"valid", user.Email, user, nil},
		{"notfound", "nosuchuser@example.com", nil, app.ErrUserNotFound},
	}

	for _, expected := range tests {
		t.Run(expected.name, func(t *testing.T) {
			u, err := service.Get(expected.email)

			if err != expected.err {
				t.Errorf("invalid error: expected %v got %v", expected.err, err)
			}

			if diff := deep.Equal(u, expected.user); diff != nil {
				t.Errorf("invalid user: %v", diff)
			}
		})
	}
}
