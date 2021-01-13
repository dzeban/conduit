package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

// GetUser returns user by email from Postgres store
func (s *Store) GetUser(email string) (*app.User, error) {
	query := `
		SELECT
			id,
			name,
			bio,
			image,
			password_hash
		FROM
			users
		WHERE
			email = $1
	`

	row := s.db.QueryRowx(query, email)

	// Scan the row using simple Scan method.
	// We can't use StructScan to the app.User var because bio and image may be
	// NULL so these fields must be handled via sql.NullString. We can't use
	// these sql-specific types in app.User because they're, well, sql-specific
	var id int
	var name, passwordHash string
	var bio, image sql.NullString
	err := row.Scan(&id, &name, &bio, &image, &passwordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	user := app.User{
		Id:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio.String,
		Image:        image.String,
	}

	return &user, nil
}

func (s *Store) GetUserById(id int) (*app.User, error) {
	query := `
		SELECT
			name,
			email,
			bio,
			image,
			password_hash
		FROM
			users
		WHERE
			id = $1
	`

	row := s.db.QueryRowx(query, id)

	// Scan the row using simple Scan method.
	// We can't use StructScan to the app.User var because bio and image may be
	// NULL so these fields must be handled via sql.NullString. We can't use
	// these sql-specific types in app.User because they're, well, sql-specific
	var name, email, passwordHash string
	var bio, image sql.NullString
	err := row.Scan(&name, &email, &bio, &image, &passwordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	user := app.User{
		Id:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          bio.String,
		Image:        image.String,
	}

	return &user, nil
}

// AddUser adds new user to the Postgres user store and returns it
func (s *Store) AddUser(user *app.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, bio, image)
		VALUES (:name, :email, :password_hash, :bio, :image)
	`

	_, err := s.db.NamedExec(query, &user)
	if err != nil {
		return errors.Wrap(err, "failed to insert user to db")
	}

	return nil
}

// UpdateUser modifies user by email and return updated user object
func (s *Store) UpdateUser(user *app.User) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err :=
		psql.Update("users").
			SetMap(user.Map()).
			Where(sq.Eq{"id": user.Id}).
			ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}

	// Execute update.
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}
