package postgres

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

func (s *Store) GetProfile(username string) (*app.Profile, error) {
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
		return nil, app.ErrorUserNotFound
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

func (s Store) FollowProfile(follower, follows string) error {
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

func (s Store) UnfollowProfile(follower, follows string) error {
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
