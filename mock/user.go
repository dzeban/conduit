package mock

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
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
				Name:  "user1",
				Email: "user1@example.com",
				//user1pass
				Password: "$argon2id$v=19$m=32768,t=5,p=1$1Eg31Vt/wwNwSvL4fIl2AA$DH1fejssBGLRpUBtMFwmaf7x7DCOJtcDsVWYvulkkZxnqKWWJCyUQgv5RZTJSSn7CzILo8cgGCFxAM8pnShZyw",
			},
			{
				Name:  "user2",
				Email: "user2@example.com",
				Bio:   "user2 bio",
				// user2pass
				Password: "$argon2id$v=19$m=32768,t=5,p=1$jsekOp1Q4F7l00w7rORgfw$mGPg4IdawxwABBdvKESFOEYr9ZZbFA92Q97KxJGR4PSFfRGMyGcsD7lTq+/LKfxkclTxWpVr7RLbrc17uZyYZw",
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

	return nil, app.ErrUserNotFound
}

// Login checks email and password and returns the user object
func (s *UserService) Login(req app.UserRequest) (*app.User, error) {
	u, err := s.Get(req.User.Email)
	if err != nil {
		return nil, app.ErrUserNotFound
	}

	ok, err := password.Check(req.User.Password, u.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check password during login")
	}

	if !ok {
		return nil, fmt.Errorf("password mismatch")
	}

	return u, nil
}

// Register creates new user in the service and returns it
func (s *UserService) Register(req app.UserRequest) (*app.User, error) {
	u, _ := s.Get(req.User.Email)
	if u != nil {
		return nil, app.ErrUserExists
	}

	hash, err := password.HashAndEncode(req.User.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create password hash")
	}

	newUser := req.User
	newUser.Password = hash

	s.users = append(s.users, newUser)
	return &newUser, nil
}

// Update overwrite user found by
// func (s *UserService) Update(email string, newData app.User) (*app.User, error) {
// 	u, err := s.Get(email)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to update user")
// 	}

// 	// Update with non-empty fields
// 	if newData.Name != "" {
// 		u.Name = newData.Name
// 	}

// 	if newData.Email != "" {
// 		u.Email = newData.Email
// 	}

// 	if newData.Password != "" {
// 		u.Password = newData.Password
// 	}

// 	if newData.Bio != "" {
// 		u.Bio = newData.Bio
// 	}

// 	if newData.Image != nil {
// 		u.Image = newData.Image
// 	}

// 	// TODO: regenerate token

// 	return u, nil
// }
