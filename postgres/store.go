package postgres

import (
	"fmt"
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

	db.MapperFunc(CamelToSnakeASCII)

	return &Store{db: db}, nil
}

// CamelToSnakeASCII converts camel case strings to snake case.
// It's used as a mapper for sqlx.
// It's a simplified version of the same name function found in
// https://github.com/scylladb/go-reflectx.
func CamelToSnakeASCII(s string) string {
	buf := []byte(s)
	out := make([]byte, 0, len(buf)+3)

	l := len(buf)
	for i := 0; i < l; i++ {
		if !(allowedChar(buf[i]) || buf[i] == '_') {
			panic(fmt.Sprint("not allowed name ", s))
		}

		b := buf[i]

		if isUpper(b) {
			if i > 0 && buf[i-1] != '_' && (isLower(buf[i-1]) || (i+1 < l && isLower(buf[i+1]))) {
				out = append(out, '_')
			}
			b = toLower(b)
		}

		out = append(out, b)
	}

	return string(out)
}

func isUpper(b byte) bool {
	return (b >= 'A' && b <= 'Z')
}

func isLower(b byte) bool {
	return !isUpper(b)
}

func toLower(b byte) byte {
	if isUpper(b) {
		return b + 32
	}
	return b
}

func allowedChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}
