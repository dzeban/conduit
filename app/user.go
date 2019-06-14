package app

import (
	"encoding/json"
	"errors"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("no such user")
)

// UserService defines and interface to work with users
type UserService interface {
	Login(ur UserRequest) (*User, error)
	Register(ur UserRequest) (*User, error)
	Get(email string) (*User, error)
	// Update(email string, newData User) (*User, error)
}

// UserServiceConfig describes configuration for UserService
type UserServiceConfig struct {
	Type   string `default:"postgres"`
	DSN    string `default:"postgres://postgres:postgres@postgres/conduit?sslmode=disable"`
	Secret string
}

// User represents a user
// User is identified by email and authenticated by JWT
// Password is hidden by custom marshaller
type User struct {
	Name     string `json:"username" db:"name"`
	Email    string `json:"email" db:"email"`
	Bio      string `json:"bio,omitempty" db:"bio"`
	Image    string `json:"image,omitempty" db:"image"` // base64 encoded
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty" db:"password"`
}

// MarshalJSON custom serializer hides password field
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string `json:"username"`
		Email string `json:"email"`
		Bio   string `json:"bio,omitempty"`
		Image string `json:"image,omitempty"`
		Token string `json:"token,omitempty"`
	}{
		Name:  u.Name,
		Email: u.Email,
		Bio:   u.Bio,
		Image: u.Image,
		Token: u.Token,
	})
}

// UserRequest represent a json structure used
// in user serice requests and responses.
type UserRequest struct {
	User User `json:"user"`
}

// MarshalJSON custom serializer hides password field
func (ur UserRequest) MarshalJSON() ([]byte, error) {
	u := ur.User

	// Marshal via anonymous struct identical to UserRequest
	// to avoid infinite recursion with marshalling.
	// This works because anonymous struct is a different type.
	return json.Marshal(&struct {
		User User `json:"user"`
	}{
		User: User{
			Name:  u.Name,
			Email: u.Email,
			Token: u.Token,
		},
	})
}

// ValidateForRegister validates user object that is used in Register handler
func (u User) ValidateForRegister() error {
	if u.Name == "" {
		return errors.New("username is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// ValidateForLogin validates user object that is used in Login handler
func (u User) ValidateForLogin() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
