package user

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

// buildUpdateUserQuery constructs update query from non-nil field of User
// It returns the prebuilt query and user args for binding
func buildUpdateUserQuery(db *sqlx.DB, user *app.User) (string, []interface{}, error) {
	// Check for empty struct
	if *user == (app.User{}) {
		return "", nil, errors.New("empty user struct")
	}

	query := "UPDATE users SET "
	if user.Name != "" {
		query += "name = :name, "
	}
	if user.Email != "" {
		query += "email = :email, "
	}
	if user.Bio != "" {
		query += "bio = :bio, "
	}
	if user.Image != "" {
		query += "image = :image, "
	}
	if user.Password != "" {
		query += "password = :password, "
	}

	// Cut last comma
	query = strings.TrimSuffix(query, ", ")

	// Convert to ? bindvars to append positional WHERE condition later
	query, args, err := sqlx.Named(query, user)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to bind struct to query")
	}

	// Add condition using positional bindvars
	query += " WHERE email = ?"

	// Rebind to Postgres bindvars ($1, $2, etc.)
	query = db.Rebind(query)

	return query, args, nil
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

func ServerError(err error, msg string) string {
	// errors.Wrap doesn't handle nil errors. To avoid nil pointer in error
	// message we create empty error here when error is nil
	if err == nil {
		err = errors.New("")
	}
	return fmt.Sprintf(`{"error":{"message":["%s"]}}`, errors.Wrap(err, msg))
}
