package user

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
)

type RegisterRequest struct {
	User RegisterUser `json:"user"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // NOTE: Plaintext password from user
}

func (r *RegisterRequest) Validate() error {
	if r.User.Username == "" {
		return errorUsernameIsRequired
	}

	if r.User.Email == "" {
		return errorEmailIsRequired
	}

	if r.User.Password == "" {
		return errorPasswordIsRequired
	}

	return nil
}

// Register creates new user in the service
func (s *Service) Register(req *RegisterRequest) (*app.User, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ServiceError(err)
	}

	// Check if user exists
	u, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get user"))
	}

	if u != nil {
		return nil, app.ServiceError(errorUserExists)
	}

	// Replace password with hash
	hash, err := password.HashAndEncode(req.User.Password)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to create password hash"))
	}

	user := &app.User{
		Name:         req.User.Username,
		Email:        req.User.Email,
		PasswordHash: hash,
	}

	// Store new user
	err = s.store.AddUser(user)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to add new user"))
	}

	// Return added user
	u, err = s.store.GetUser(user.Email)
	if err != nil {
		return nil, app.InternalError(errorUserNotCreated)
	}

	return u, nil
}
