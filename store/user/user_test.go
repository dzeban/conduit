package user

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/jmoiron/sqlx"

	"github.com/dzeban/conduit/app"
)

func TestBuildUpdateUserQuery(t *testing.T) {
	tests := []struct {
		user   *app.User
		query  string
		args   []interface{}
		hasErr bool
	}{
		{
			user:   &app.User{},
			query:  "",
			args:   nil,
			hasErr: true,
		},
		{
			user:   &app.User{Name: "test"},
			query:  "UPDATE users SET name = $1 WHERE email = $2",
			args:   []interface{}{"test"},
			hasErr: false,
		},
		{
			user:   &app.User{Email: "test@example.com", Bio: "mystery"},
			query:  "UPDATE users SET email = $1, bio = $2 WHERE email = $3",
			args:   []interface{}{"test@example.com", "mystery"},
			hasErr: false,
		},
	}

	db, err := sqlx.Open("postgres", "")
	if err != nil {
		t.Error(err)
	}

	for _, test := range tests {
		query, args, err := buildUpdateUserQuery(db, test.user)

		if test.query != query {
			t.Errorf("query doesn't match, expected '%v', got '%v'\n", test.query, query)
		}

		if diff := deep.Equal(args, test.args); diff != nil {
			t.Errorf("invalid args: %v", diff)
		}

		hasErr := (err != nil)
		if hasErr != test.hasErr {
			t.Errorf("err doesn't match, expected '%v', got '%v'\n", test.hasErr, hasErr)
		}
	}
}
