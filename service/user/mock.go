package user

import (
	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
	"github.com/pkg/errors"
)

var testService = New(newMockStore(), "test")

type mockStore struct {
	users     map[string]app.User
	followers map[string][]string
}

func newMockStore() app.UserStore {
	hash, _ := password.HashAndEncode("test")
	users := make(map[string]app.User)
	users["test@example.com"] = app.User{
		Name:     "test",
		Email:    "test@example.com",
		Password: hash,
	}

	followers := make(map[string][]string)
	followers["test_follower"] = make([]string, 1)
	followers["test_follower"] = append(followers["test_follower"], "test_follow")
	return mockStore{users, followers}
}

func (s mockStore) Get(email string) (*app.User, error) {
	u, ok := s.users[email]
	if !ok {
		return nil, app.ErrUserNotFound
	}

	return &u, nil
}

func (s mockStore) Add(user app.User) error {
	s.users[user.Email] = user
	return nil
}

func (s mockStore) Update(email string, user app.User) error {
	s.users[email] = user
	return nil
}

func (s mockStore) Profile(username string) (*app.Profile, error) {
	var user *app.User
	for _, u := range s.users {
		if u.Name == username {
			user = &u
		}
	}

	if user == nil {
		return nil, app.ErrUserNotFound
	}

	return &app.Profile{
		Name: user.Name,
	}, nil
}

func (s mockStore) Follow(follower, following string) error {
	_, ok := s.followers[follower]
	if !ok {
		s.followers[follower] = make([]string, 1)
	}

	s.followers[follower] = append(s.followers[follower], following)

	return nil
}

func (s mockStore) Unfollow(follower, following string) error {
	_, ok := s.followers[follower]
	if !ok {
		return errors.New("no follow relationship")
	}

	f := s.followers[follower]
	for i, v := range f {
		if v == following {
			f = append(f[:i], f[i+1:]...)
		}
	}

	return nil
}
