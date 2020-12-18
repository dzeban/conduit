package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/db"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(DSN string) (*Store, error) {
	db, err := db.ConnectLoop("postgres", DSN, 1*time.Minute)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to users db")
	}

	return &Store{db: db}, nil
}
