package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/transport"
)

var (
	errorJWTNoAuthorizationHeader = errors.New("no Authorization header")
)

type AuthType int

const (
	AuthTypeRequired AuthType = iota
	AuthTypeOptional
)

func Auth(secret []byte, typ AuthType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return transport.WithError(func(w http.ResponseWriter, r *http.Request) error {
			authHeader, ok := r.Header["Authorization"]
			if !ok {
				if typ == AuthTypeRequired {
					return app.AuthError(errorJWTNoAuthorizationHeader)
				} else {
					// If authorization header is absent pass to the next handler
					// without updating context.
					next.ServeHTTP(w, r)
					return nil
				}
			}

			u, err := userFromJWT(authHeader[0], secret)
			if err != nil {
				return app.AuthError(errors.Wrap(err, "invalid JWT"))
			}

			// Store user in context for the further reference
			authCtx := u.NewContext(r.Context())

			next.ServeHTTP(w, r.WithContext(authCtx))
			return nil
		})
	}
}

// userFromJWT takes JWT from Authorization header, parses and
// validates it and returns user. JWT is expected in "Token <token>" format.
func userFromJWT(header string, secret []byte) (*app.User, error) {
	tokenVals := strings.Split(header, " ")

	if len(tokenVals) != 2 {
		return nil, errors.New("invalid auth header format, expected 2 elements")
	}

	if tokenVals[0] != "Token" {
		return nil, fmt.Errorf("invalid auth header format, expected Token <token>, got %#v", header)
	}

	claims, err := Parse(tokenVals[1], secret)
	if err != nil {
		return nil, errors.Wrap(err, "jwt parsing")
	}

	return claims.User, nil
}
