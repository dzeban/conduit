package user

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
)

// LoginRequest describes request JSON for login handler
type LoginRequest struct {
	User LoginUser `json:"user"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"` // NOTE: Plaintext password from user
}

func (r *LoginRequest) Validate() error {
	if r.User.Email == "" {
		return errorEmailIsRequired
	}

	if r.User.Password == "" {
		return errorPasswordIsRequired
	}

	return nil
}

// Login checks the user request and logins the user
func (s *Service) Login(req *LoginRequest) (*app.User, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ServiceError(err)
	}

	// Lookup user by email
	user, err := s.store.GetUser(req.User.Email)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get user"))
	}

	if user == nil {
		return nil, app.ServiceError(errorUserNotFound)
	}

	// Check password
	ok, err := password.Check(req.User.Password, user.PasswordHash)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to check password during login"))
	}

	if !ok {
		return nil, app.AuthError(errorPasswordMismatch)
	}

	// Return the user
	return user, nil
}
