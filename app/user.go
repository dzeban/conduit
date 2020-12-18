package app

import (
	"context"
)

// User represents a user
// User is identified by email and authenticated by JWT.
// Password is hidden by custom marshaller.
type User struct {
	Name  string `json:"username"`
	Email string `json:"email"`
	Bio   string `json:"bio,omitempty"`

	// Image is base64-encoded user avatar
	Image string `json:"image,omitempty"`

	// Token is JWT returned by user register, login and update
	Token string `json:"token,omitempty"`

	PasswordHash string `json:"-"`
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
	if u.PasswordHash != "" {
		m["password_hash"] = u.PasswordHash
	}

	return m
}

// context.Context helpers
type contextKey string

func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey(""), u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(contextKey("")).(*User)
	return u, ok
}
