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
		return app.ErrorValidationUsernameIsRequired
	}

	if r.User.Email == "" {
		return app.ErrorValidationEmailIsRequired
	}

	if r.User.Password == "" {
		return app.ErrorValidationPasswordIsRequired
	}

	return nil
}

// Register creates new user in the service
func (s *Service) Register(req *RegisterRequest) (*app.User, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ValidationError(err)
	}

	// Check if user exists
	u, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get user"))
	}

	if u != nil {
		return nil, app.ServiceError(app.ErrorUserExists)
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

	return user, nil
}
