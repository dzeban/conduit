package user

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/db"
)

type PostgresStore struct {
	db *sqlx.DB
}

func New(DSN string) (app.UserStore, error) {
	db, err := db.ConnectLoop("postgres", DSN, 1*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to users db")
	}

	return PostgresStore{db: db}, nil
}

// Get returns user by email from Postgres store
func (s PostgresStore) Get(email string) (*app.User, error) {
	query := `
		SELECT
			name,
			bio,
			image,
			password
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
	var name, password string
	var bio, image sql.NullString
	err := row.Scan(&name, &bio, &image, &password)
	if err == sql.ErrNoRows {
		return nil, app.ErrUserNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query user")
	}

	user := app.User{
		Name:     name,
		Email:    email,
		Password: password,
		Bio:      bio.String,
		Image:    image.String,
	}

	return &user, nil
}

// Add adds new user to the Postgres user store and returns it
func (s PostgresStore) Add(user app.User) error {
	query := `
		INSERT INTO users (name, email, password, bio, image)
		VALUES (:name, :email, :password, :bio, :image)
	`

	_, err := s.db.NamedExec(query, &user)
	if err != nil {
		return errors.Wrap(err, "failed to insert user to db")
	}

	return nil
}

// Update modifies user by email and return updated user object
func (s PostgresStore) Update(email string, user app.User) error {
	query, args, err := buildUpdateUserQuery(s.db, &user)
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}

	// Execute update.
	// args are created from user struct and used in SET expressions for actual
	// update. email is used in WHERE clause as search condition.
	args = append(args, email)
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute update query")
	}

	return nil
}

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

func (s PostgresStore) Profile(username string) (*app.Profile, error) {
	query := `
		SELECT
			name,
			bio,
			image
		FROM
			users
		WHERE
			name = $1
	`

	row := s.db.QueryRowx(query, username)

	// Scan the row using simple Scan method.
	// We can't use StructScan to the app.User var because bio and image may be
	// NULL so these fields must be handled via sql.NullString. We can't use
	// these sql-specific types in app.User because they're, well, sql-specific
	var name string
	var bio, image sql.NullString
	err := row.Scan(&name, &bio, &image)
	if err == sql.ErrNoRows {
		return nil, app.ErrUserNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query profile")
	}

	profile := app.Profile{
		Name:  name,
		Bio:   bio.String,
		Image: image.String,
	}

	return &profile, nil
}

func (s PostgresStore) Follow(follower, follows string) error {
	query := `
		INSERT INTO followers (follower, follows)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := s.db.Exec(query, follower, follows)
	if err != nil {
		return errors.Wrap(err, "failed to add follow relationship to db")
	}

	return nil
}

func (s PostgresStore) Unfollow(follower, follows string) error {
	query := `
		DELETE FROM followers
		WHERE follower=$1 AND follows=$2
	`

	_, err := s.db.Exec(query, follower, follows)
	if err != nil {
		return errors.Wrap(err, "failed to delete follow relationship from db")
	}

	return nil
}