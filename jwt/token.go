package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func New(email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    email,
		"signed": true,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", errors.Wrap(err, "failed to create token")
	}

	return tokenString, nil
}
