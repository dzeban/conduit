package user

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

var (
	errorEmailIsRequired     = errors.New("email is required")
	errorPasswordIsRequired  = errors.New("password is required")
	errorPasswordMismatch    = errors.New("password mismatch")
	errorUsernameIsRequired  = errors.New("username is required")
	errorUserExists          = errors.New("user exists")
	errorUserNotFound        = errors.New("user not found")
	errorUserUpdateForbidden = errors.New("user update forbidden")
	errorInvalidRequest      = errors.New("invalid request")
	errorUserNotCreated      = errors.New("user not created")
)

type Store interface {
	GetUser(email string) (*app.User, error)
	GetUserById(id int) (*app.User, error)
	AddUser(user *app.User) error
	UpdateUser(user *app.User) error
}

// Service provides a service for interacting with user accounts
type Service struct {
	store Store
}

// NewService creates new instance of the service with provided store
func NewService(store Store) *Service {
	return &Service{store}
}

// Get returns user by email
func (s *Service) Get(email string) (*app.User, error) {
	u, err := s.store.GetUser(email)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get user"))
	}

	if u == nil {
		return nil, app.ServiceError(errorUserNotFound)
	}

	return u, nil
}
