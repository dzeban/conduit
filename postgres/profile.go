package postgres

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

// PostgresProfile is the same as app.Profile but suitable for StructScan
// method because it has sql.Null* types.
type PostgresProfile struct {
	Id        int
	Name      string
	Bio       sql.NullString
	Image     sql.NullString
	Following bool
}

func (s *Store) GetProfile(username string, follower *app.Profile) (*app.Profile, error) {
	query := `
		SELECT
			id,
			name,
			bio,
			image,
			followee IS NOT NULL AS following
		FROM users u
		LEFT JOIN followers f
		ON (u.id = f.followee AND f.follower = $1)
		WHERE u.name = $2;
	`

	// Set follower id if it's a request for authenticated user.
	// If follower is not set then id will be 0 and "following" will always be
	// False because ids start with 1.
	followerId := 0
	if follower != nil {
		followerId = follower.Id
	}

	row := s.db.QueryRowx(query, followerId, username)

	var p PostgresProfile
	err := row.StructScan(&p)
	if err == sql.ErrNoRows {
		return nil, app.ErrorProfileNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to query profile")
	}

	profile := app.Profile{
		Id:        p.Id,
		Name:      p.Name,
		Bio:       p.Bio.String,
		Image:     p.Image.String,
		Following: p.Following,
	}

	return &profile, nil
}

func (s Store) FollowProfile(follower, followee *app.Profile) error {
	query := `
		INSERT INTO followers (follower, followee)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	_, err := s.db.Exec(query, follower.Id, followee.Id)
	if err != nil {
		return errors.Wrap(err, "failed to add follow relationship to db")
	}

	return nil
}

func (s Store) UnfollowProfile(follower, followee *app.Profile) error {
	query := `
		DELETE FROM followers
		WHERE follower = $1 AND followee = $2
	`

	_, err := s.db.Exec(query, follower.Id, followee.Id)
	if err != nil {
		return errors.Wrap(err, "failed to delete follow relationship from db")
	}

	return nil
}
