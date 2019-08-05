package user

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		user app.User
		err  error
	}{
		{
			app.User{Email: "test@example.com", Name: "test", Password: "test"},
			nil,
		},
		{
			app.User{Email: "no@example.com", Password: "test"},
			app.ErrUserNotFound,
		},
		{
			app.User{Email: "test@example.com", Password: "invalid"},
			app.ErrPasswordMismatch,
		},
	}

	s := New(newMockStore(), "test")
	for _, c := range cases {
		u, err := s.Login(c.user)
		if err != c.err {
			t.Errorf("Login(%#v) => unexpected error, want %v, got %v", c.user, c.err, err)
		}

		if err == nil {
			checkUser(c.user, *u, t)
		}
	}
}

func TestRegister(t *testing.T) {
	cases := []struct {
		user app.User
		err  error
	}{
		{
			app.User{Email: "test@example.com", Password: "test"},
			app.ErrUserExists,
		},
		{
			app.User{Email: "new@example.com", Password: "new"},
			nil,
		},
	}

	s := New(newMockStore(), testSecret)
	for _, c := range cases {
		u, err := s.Register(c.user)
		if err != c.err {
			t.Errorf("Register(%#v) => unexpected error, want %v, got %v", c.user, c.err, err)
		}

		if err == nil {
			checkUser(c.user, *u, t)
		}
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		user   app.User
		update app.User
		want   app.User
	}{
		{
			app.User{Email: "passwordUpdate@example.com", Password: "old"},
			app.User{Password: "new"},
			app.User{Email: "passwordUpdate@example.com", Password: "new"},
		},
		{
			app.User{Email: "fieldsUpdate@example.com", Password: "password"},
			app.User{Name: "name", Bio: "bio", Image: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+ip1sAAAAASUVORK5CYII="},
			app.User{Email: "fieldsUpdate@example.com", Password: "password", Name: "name", Bio: "bio", Image: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+ip1sAAAAASUVORK5CYII="},
		},
	}

	s := New(newMockStore(), testSecret)
	for _, c := range cases {
		_, err := s.Register(c.user)
		if err != nil {
			t.Errorf("Register(%#v) => unexpected error %v", c.user, err)
		}

		u, err := s.Update(c.user.Email, c.update)
		if err != nil {
			t.Errorf("Update(%v, %#v) => unexpected error %v", c.user.Email, c.user, err)
		}

		checkUser(c.want, *u, t)
	}
}

func checkUser(want, got app.User, t *testing.T) {
	// First, check passwords
	// want has plaintext password, got has hashed password
	ok, err := password.Check(want.Password, got.Password)
	if err != nil {
		t.Errorf("failed to check password: %s", err)
	}

	if !ok {
		t.Errorf("password mismatch")
	}

	// Now, remove password and deep check structs
	want.Password = ""
	got.Password = ""

	if diff := deep.Equal(want, got); diff != nil {
		t.Errorf("users don't match: %s", diff)
	}
}
