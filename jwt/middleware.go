package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dzeban/conduit/app"
	"github.com/pkg/errors"
)

var (
	ErrJWTNoAuthorizationHeader = errors.New("no Authorization header")
)

func Auth(secret []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader, ok := r.Header["Authorization"]
			if !ok {
				http.Error(w, app.ServerError(ErrJWTNoAuthorizationHeader, ""), http.StatusUnauthorized)
				return
			}

			user, err := userFromJWT(authHeader[0], secret)
			if err != nil {
				http.Error(w, app.ServerError(err, "failed to parse jwt"), http.StatusUnauthorized)
				return
			}

			// Store user in context for the further reference
			authCtx := app.UserNewContext(r.Context(), user)

			next.ServeHTTP(w, r.WithContext(authCtx))
		}

		return http.HandlerFunc(fn)
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
