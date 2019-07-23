package user

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// parseJWTClaimsFromHeader takes JWT from Authorization header, parses and
// validates it and returns claims. JWT is expected in "Token <token>" format.
func parseJWTClaimsFromHeader(header string, secret []byte) (map[string]interface{}, error) {
	tokenVals := strings.Split(header, " ")

	if len(tokenVals) != 2 {
		return nil, errors.New("invalid auth header format, expected 2 elements")
	}

	if tokenVals[0] != "Token" {
		return nil, fmt.Errorf("invalid auth header format, expected Token <token>, got %#v", header)
	}

	token, err := jwt.Parse(tokenVals[1], func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt parsing error")
	}

	if !token.Valid {
		return nil, errors.New("jwt is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to cast claims to jwt.MapClaims")
	}

	return claims, nil
}

func ServerError(err error, msg string) string {
	// errors.Wrap doesn't handle nil errors. To avoid nil pointer in error
	// message we create empty error here when error is nil
	if err == nil {
		err = errors.New("")
	}
	return fmt.Sprintf(`{"error":{"message":["%s"]}}`, errors.Wrap(err, msg))
}
