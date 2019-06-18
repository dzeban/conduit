package user

import (
	"strings"

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
