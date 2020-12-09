package app

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("no such user")
	ErrPasswordMismatch = errors.New("password mismatch")
)

type UserStore interface {
	Get(email string) (*User, error)
	Add(user User) error
	Update(email string, user User) error
	Profile(username string) (*Profile, error)
	Follow(follower, follows string) error
	Unfollow(follower, follows string) error
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
	Name  string `json:"username"`
	Email string `json:"email"`
	Bio   string `json:"bio,omitempty"`

	// Image is base64-encoded user avatar
	Image string `json:"image,omitempty"`

	// Token is JWT returned by user register, login and update
	Token string `json:"token,omitempty"`

	// Password is used to carry plain-text password from user login, register,
	// update requests and then immediately replaced with hash.
	// This multi-purpose'ness appears because User struct is shared for request
	// and for app type.
	//
	// TODO: Introduce struct for Login/Register/Update requests, rename field
	// here to hash. This also will make custom marshaller obsolete.
	Password string `json:"password,omitempty"`
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

func (u User) Map() map[string]interface{} {
	m := make(map[string]interface{})

	if u.Name != "" {
		m["name"] = u.Name
	}
	if u.Email != "" {
		m["email"] = u.Email
	}
	if u.Bio != "" {
		m["bio"] = u.Bio
	}
	if u.Image != "" {
		m["image"] = u.Image
	}
	if u.Password != "" {
		m["password"] = u.Password
	}

	return m
}

// context.Context helpers
type userContextKey string

func UserNewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userContextKey(""), u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey("")).(*User)
	return u, ok
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

// Profile is a public user info with restricted set of fields
type Profile struct {
	Name      string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`     // base64 encoded
	Following bool   `json:"following"` // set for authenticated users
}

// ProfileResponse is a structure returned in profile handlers
type ProfileResponse struct {
	Profile Profile `json:"profile"`
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

// ValidateForUpdate validates user object that is used in Update handler
// It checks the presense of at least one field.
func (u User) ValidateForUpdate() error {
	if u.Name == "" &&
		u.Email == "" &&
		u.Bio == "" &&
		u.Image == "" &&
		u.Password == "" {
		return errors.New("at least one of name, email, bio, image, password is required for update")
	}

	return nil
}
