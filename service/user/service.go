package user

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
	"github.com/dzeban/conduit/store/user"
)

// Service provides a service for interacting with user accounts
type Service struct {
	store  app.UserStore
	router *chi.Mux
	secret []byte
}

// New is a constructor for a Service
func New(store app.UserStore, secret string) *Service {
	router := chi.NewRouter()

	s := &Service{
		store:  store,
		router: router,
		secret: []byte(secret),
	}

	router.Post("/", s.HandleUserRegister)
	router.Post("/login", s.HandleUserLogin)
	router.Get("/", s.jwtAuthHandler(s.HandleUserGet))
	router.Put("/", s.jwtAuthHandler(s.HandleUserUpdate))

	router.Get("/{username}", s.HandleProfileGet)
	router.Post("/{username}/follow", s.jwtAuthHandler(s.HandleProfileFollow))
	router.Post("/{username}/unfollow", s.jwtAuthHandler(s.HandleProfileUnfollow))

	return s
}

// NewFromDSN wraps New to create user service from DSN instead of explicit stores
func NewFromDSN(DSN, secret string) (*Service, error) {
	store, err := user.New(DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create user store for DSN %s", DSN)
	}

	return New(store, secret), nil
}

// Get returns user by email
func (s *Service) Get(email string) (*app.User, error) {
	return s.store.Get(email)
}

// Login checks the user request and logins the user
func (s *Service) Login(user app.User) (*app.User, error) {
	u, err := s.store.Get(user.Email)
	if err != nil {
		return nil, app.ErrUserNotFound
	}

	ok, err := password.Check(user.Password, u.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check password during login")
	}

	if !ok {
		return nil, app.ErrPasswordMismatch
	}

	return u, nil
}

// Update modifies user found by email with the new data passed in user.
// It returns updated user.
func (s *Service) Update(email string, user app.User) (*app.User, error) {
	// If password is being changed, make the hash from it
	if user.Password != "" {
		hash, err := password.HashAndEncode(user.Password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create password hsah")
		}

		user.Password = hash
	}

	err := s.store.Update(email, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	// Return updated user
	return s.store.Get(email)
}

func (s *Service) Register(user app.User) (*app.User, error) {
	// Check if user exists
	u, _ := s.store.Get(user.Email)
	if u != nil {
		return nil, app.ErrUserExists
	}

	// Replace password with hash
	hash, err := password.HashAndEncode(user.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create password hsah")
	}

	newUser := user
	newUser.Password = hash

	err = s.store.Add(newUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new user")
	}

	return &newUser, nil
}
