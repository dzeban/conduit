package user

import (
	"database/sql"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
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
	db, err := sqlx.Connect("postgres", DSN)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to users db")
	}

	router := mux.NewRouter().StrictSlash(true)

	s := &Service{db: db, router: router, secret: []byte(secret)}

	router.HandleFunc("/users/", s.HandleUserRegister).Methods("POST")
	router.HandleFunc("/users/login", s.HandleUserLogin).Methods("POST")
	router.HandleFunc("/users/", s.jwtAuthHandler(s.HandleUserGet)).Methods("GET")

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
	// NULL so these fields must be handled via sql.NullString.
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

// Login checks the user request and logins the user
func (s *Service) Login(req app.UserRequest) (*app.User, error) {
	fmt.Println(req.User.Email)
	u, err := s.Get(req.User.Email)
	if err != nil {
		return nil, app.ErrUserNotFound
	}

	fmt.Println(req.User.Password, u.Password)
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