package user

import (
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/jwt"
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

	// Unauthenticated endpoints
	router.Post("/", s.HandleUserRegister)
	router.Post("/login", s.HandleUserLogin)
	router.Get("/{username}", s.HandleProfileGet)

	// Endpoints protected by JWT auth
	router.Group(func(r chi.Router) {
		r.Use(jwt.Auth(s.secret))

		r.Get("/", s.HandleUserGet)
		r.Put("/", s.HandleUserUpdate)
		r.Post("/{username}/follow", s.HandleProfileFollow)
		r.Post("/{username}/unfollow", s.HandleProfileUnfollow)
	})

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
	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	if user.Password == "" {
		return nil, errors.New("password is required")
	}

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
	if user.Name == "" &&
		user.Email == "" &&
		user.Bio == "" &&
		user.Image == "" &&
		user.Password == "" {
		return nil, errors.New("at least one of name, email, bio, image, password is required for update")
	}

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
	if user.Name == "" {
		return nil, errors.New("username is required")
	}

	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	if user.Password == "" {
		return nil, errors.New("password is required")
	}

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

func (s *Service) GetProfile(username string) (*app.Profile, error) {
	return s.store.Profile(username)
}

func (s *Service) FollowProfile(follower *app.User, username string) (*app.Profile, error) {
	// Query profile to follow to ensure it exists
	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		return nil, errors.New("no such profile")
	} else if err != nil {
		return nil, errors.New("failed to get profile")
	}

	err = s.store.Follow(follower.Name, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to follow user")
	}

	return profile, nil
}

func (s *Service) UnfollowProfile(follower *app.User, username string) (*app.Profile, error) {
	// Query profile to follow to ensure it exists
	profile, err := s.store.Profile(username)
	if err == app.ErrUserNotFound {
		return nil, errors.New("no such profile")
	} else if err != nil {
		return nil, errors.New("failed to get profile")
	}

	err = s.store.Unfollow(follower.Name, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unfollow user")
	}

	return profile, nil
}
