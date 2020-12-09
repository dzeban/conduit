package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

type TokenClaims struct {
	User *app.User `json:"user"`
}

func (t TokenClaims) Valid() error {
	if t.User == nil {
		return errors.New("empty user")
	}

	if t.User.Name == "" {
		return errors.New("username is empty")
	}

	if t.User.Email == "" {
		return errors.New("email is empty")
	}

	return nil
}

func New(user *app.User, secret []byte) (string, error) {
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, TokenClaims{user}).
		SignedString(secret)
}

func Parse(token string, secret []byte) (*TokenClaims, error) {
	var tc TokenClaims
	_, err := jwt.ParseWithClaims(token, &tc, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "jwt parsing error")
	}

	return &tc, nil
}
