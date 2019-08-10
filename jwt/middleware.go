package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dzeban/conduit/app"
	"github.com/pkg/errors"
)

var (
	ErrJWTNoAuthorizationHeader = errors.New("no Authorization header")
	ErrJWTNoSignedClaim         = errors.New("token does not have signed claim")
	ErrJWTNoSubClaim            = errors.New("no sub claim")
)

func Auth(secret []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader, ok := r.Header["Authorization"]
			if !ok {
				http.Error(w, app.ServerError(ErrJWTNoAuthorizationHeader, ""), http.StatusUnauthorized)
				return
			}

			claims, err := parseJWTClaimsFromHeader(authHeader[0], secret)
			if err != nil {
				http.Error(w, app.ServerError(err, "failed to parse jwt"), http.StatusBadRequest)
				return
			}

			if claims["signed"] != true {
				http.Error(w, app.ServerError(ErrJWTNoSignedClaim, ""), http.StatusUnauthorized)
				return
			}

			var sub interface{}
			if sub, ok = claims["sub"]; !ok {
				http.Error(w, app.ServerError(ErrJWTNoSubClaim, ""), http.StatusUnauthorized)
				return
			}

			// Store auth subject (email) to the context
			authCtx := context.WithValue(r.Context(), "email", sub)

			next.ServeHTTP(w, r.WithContext(authCtx))
		}

		return http.HandlerFunc(fn)
	}
}

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
