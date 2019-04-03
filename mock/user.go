package mock

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
)

type mockUser struct {
	user app.User
	hash string
}

// UserService implements mock user service that serve users from memory
type UserService struct {
	users []mockUser
}

// NewUserService returns new mock user service
func NewUserService() *UserService {
	return &UserService{
		users: []mockUser{
			{
				user: app.User{
					Name:  "user1",
					Email: "user1@example.com",
				},
				//user1pass
				hash: "$argon2id$v=19$m=32768,t=5,p=1$1Eg31Vt/wwNwSvL4fIl2AA$DH1fejssBGLRpUBtMFwmaf7x7DCOJtcDsVWYvulkkZxnqKWWJCyUQgv5RZTJSSn7CzILo8cgGCFxAM8pnShZyw",
			},
			{
				user: app.User{
					Name:  "user2",
					Email: "user2@example.com",
					Bio:   "user2 bio",
				},
				// user2pass
				hash: "$argon2id$v=19$m=32768,t=5,p=1$jsekOp1Q4F7l00w7rORgfw$mGPg4IdawxwABBdvKESFOEYr9ZZbFA92Q97KxJGR4PSFfRGMyGcsD7lTq+/LKfxkclTxWpVr7RLbrc17uZyYZw",
			},
		},
	}
}

// Get returns user by email
func (s *UserService) Get(email string) (*app.User, error) {
	for _, u := range s.users {
		if u.user.Email == email {
			return &u.user, nil
		}
	}

	return nil, fmt.Errorf("no user with email %s", email)
}

// // Login checks email and password and returns the user object
// func (s *UserService) Login(email, password string) (*app.User, error) {
// 	u, err := s.Get(email)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to login user")
// 	}

// 	if u.Password != password {
// 		return nil, errors.New("invalid password")
// 	}

// 	// TODO: generate token

// 	return u, nil
// }

// Register creates new user in the service and returns it
func (s *UserService) Register(user app.User, plaintextPassword string) error {
	u, _ := s.Get(user.Email)
	if u != nil {
		return fmt.Errorf("user exists")
	}

	hash, err := password.HashAndEncode(plaintextPassword)
	if err != nil {
		return errors.Wrap(err, "failed to create password hash")
	}

	mu := mockUser{
		user: user,
		hash: hash,
	}

	s.users = append(s.users, mu)
	return nil
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
