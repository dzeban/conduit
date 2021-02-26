package app

import (
	"context"
)

// User represents a user
// User is identified by email and authenticated by JWT.
// Password is hidden by custom marshaller.
type User struct {
	Id    int    `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
	Bio   string `json:"bio,omitempty"`

	// Image is base64-encoded user avatar
	Image string `json:"image,omitempty"`

	// Token is JWT returned by user register, login and update
	Token string `json:"token,omitempty"`

	PasswordHash string `json:"-"`
}

// UpdateMap returns map of fields to be updated. It sets only subset of fields
// that are allowed to be updated and set it if it contains non empty values.
func (u User) UpdateMap() map[string]interface{} {
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
type key int

var contextKey key

func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, u)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	v := ctx.Value(contextKey)
	if v == nil {
		return nil, false
	}

	u, ok := v.(*User)
	return u, ok
}
