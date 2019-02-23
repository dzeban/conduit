package mock

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

// UserService implements mock user service that serve users from memory
type UserService struct {
	users []app.User
}

// NewUserService returns new mock user service
func NewUserService() *UserService {
	return &UserService{
		users: []app.User{
			{
				Name:     "user1",
				Email:    "user1@example.com",
				Password: "user1pass",
			},
			{
				Name:     "user2",
				Email:    "user2@example.com",
				Password: "user2pass",
				Bio:      "user2 bio",
			},
		},
	}
}

// Get returns user by email
func (s *UserService) Get(email string) (*app.User, error) {
	for _, u := range s.users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("no user with email %s", email)
}

// Login checks email and password and returns the user object
func (s *UserService) Login(email, password string) (*app.User, error) {
	u, err := s.Get(email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login user")
	}

	if u.Password != password {
		return nil, errors.New("invalid password")
	}

	// TODO: generate token

	return u, nil
}

// Register creates new user in the service and returns it
func (s *UserService) Register(u app.User) (*app.User, error) {
	s.users = append(s.users, u)
	// TODO: generate token

	return &u, nil
}

// Update overwrite user found by
func (s *UserService) Update(email string, newData app.User) (*app.User, error) {
	u, err := s.Get(email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	// Update with non-empty fields
	if newData.Name != "" {
		u.Name = newData.Name
	}

	if newData.Email != "" {
		u.Email = newData.Email
	}

	if newData.Password != "" {
		u.Password = newData.Password
	}

	if newData.Bio != "" {
		u.Bio = newData.Bio
	}

	if newData.Image != nil {
		u.Image = newData.Image
	}

	// TODO: regenerate token

	return u, nil
}
