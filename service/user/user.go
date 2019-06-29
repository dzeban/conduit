package user

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/db"
	"github.com/dzeban/conduit/password"
)

// Service provides a service for interacting with user accounts
type Service struct {
	db     *sqlx.DB
	router *mux.Router
	secret []byte
}

// NewService is a constructor for a Service
func NewService(DSN string, secret string) (*Service, error) {
	db, err := db.ConnectLoop("postgres", DSN, 1*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to users db")
	}

	router := mux.NewRouter().StrictSlash(true)

	s := &Service{db: db, router: router, secret: []byte(secret)}
	router.HandleFunc("/users/", s.HandleUserRegister).Methods("POST")
	router.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
	router.HandleFunc("/users/", s.jwtAuthHandler(s.HandleUserGet)).Methods("GET")
	router.HandleFunc("/users/", s.jwtAuthHandler(s.HandleUserUpdate)).Methods("PUT")

	router.HandleFunc("/profiles/{username}", s.HandleProfileGet).Methods("GET")
	router.HandleFunc("/profiles/{username}/follow", s.jwtAuthHandler(s.HandleProfileFollow)).Methods("POST")

	return s, nil
}

// Get returns user by email
func (s *Service) Get(email string) (*app.User, error) {
	queryUser := `
		SELECT
			name,
			bio,
			image,
			password
		FROM
			users
		WHERE
			email = $1
	`

	row := s.db.QueryRowx(queryUser, email)

	// Scan the row using simple Scan method.
	// We can't use StructScan to the app.User var because bio and image may be
	// NULL so these fields must be handled via sql.NullString. We can't use
	// these sql-specific types in app.User because they're, well, sql-specific
	var name, password string
	var bio, image sql.NullString
	err := row.Scan(&name, &bio, &image, &password)
	if err == sql.ErrNoRows {
		return nil, app.ErrUserNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	user := app.User{
		Name:     name,
		Email:    email,
		Password: password,
		Bio:      bio.String,
		Image:    image.String,
	}

	return &user, nil
}

// Update modifies user by email and return updated user object
func (s *Service) Update(email string, req app.UserRequest) (*app.User, error) {
	// If password is being changed, make the hash from it
	if req.User.Password != "" {
		hash, err := password.HashAndEncode(req.User.Password)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create password hsah")
		}

		req.User.Password = hash
	}

	query, args, err := buildUpdateUserQuery(s.db, &req.User)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build update query")
	}

	// Execute update.
	// args are created from user struct and used in SET expressions for actual
	// update. email is used in WHERE clause as search condition.
	args = append(args, email)
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute update query")
	}

	// Return update user
	u, err := s.Get(email)
	if err != nil {
		return nil, app.ErrUserNotFound
	}

	return u, nil
}

// Login checks the user request and logins the user
func (s *Service) Login(req app.UserRequest) (*app.User, error) {
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
func (s *Service) Register(req app.UserRequest) (*app.User, error) {
	u, _ := s.Get(req.User.Email)
	if u != nil {
		return nil, app.ErrUserExists
	}

	hash, err := password.HashAndEncode(req.User.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create password hsah")
	}

	newUser := req.User
	newUser.Password = hash

	queryRegister := `
		INSERT INTO users (name, email, password, bio, image)
		VALUES (:name, :email, :password, :bio, :image)
	`

	_, err = s.db.NamedExec(queryRegister, &newUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert user to db")
	}

	return &newUser, nil
}

func (s *Service) Profile(username string) (*app.Profile, error) {
	queryProfile := `
		SELECT
			name,
			bio,
			image
		FROM
			users
		WHERE
			name = $1
	`

	row := s.db.QueryRowx(queryProfile, username)

	// Scan the row using simple Scan method.
	// We can't use StructScan to the app.User var because bio and image may be
	// NULL so these fields must be handled via sql.NullString. We can't use
	// these sql-specific types in app.User because they're, well, sql-specific
	var name string
	var bio, image sql.NullString
	err := row.Scan(&name, &bio, &image)
	if err == sql.ErrNoRows {
		return nil, app.ErrUserNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	profile := app.Profile{
		Name:  name,
		Bio:   bio.String,
		Image: image.String,
	}

	return &profile, nil
}

func (s *Service) Follow(follower, follows string) error {
	queryFollow := `
		INSERT INTO followers (follower, follows)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := s.db.Exec(queryFollow, follower, follows)
	if err != nil {
		return errors.Wrap(err, "failed to insert user to db")
	}

	return nil
}
