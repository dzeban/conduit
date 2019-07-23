package user

import (
	"testing"

	"github.com/dzeban/conduit/app"
)

func TestLogin(t *testing.T) {
	cases := []struct {
		user app.User
		err  error
	}{
		{
			app.User{Email: "test@example.com", Password: "test"},
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
		_, err := s.Login(c.user)
		if err != c.err {
			t.Errorf("Login(%#v) => unexpected error, want %v, got %v", c.user, c.err, err)
		}
	}
}
