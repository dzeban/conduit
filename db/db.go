package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ConnectLoop tries to connect to the DB under given DSN using a give driver
// in a loop until connection succeeds. timeout specifies the timeout for the
// loop.
func ConnectLoop(driver, DSN string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(1 * time.Second)
	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := sqlx.Connect("postgres", DSN)
			if err == nil {
				return db, nil
			}
			log.Println(errors.Wrapf(err, "failed to connect to db %s", DSN))
		}
	}
}
