package app

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUserNoPassword(t *testing.T) {
	u := User{
		Name:     "Mr. McDuck",
		Email:    "scrooge@mcduck.com",
		Bio:      "mystery...",
		Password: "superpassword",
	}

	b, err := json.Marshal(u)
	if err != nil {
		t.Error("failed to marshal User json", err)
	}

	s := string(b)

	if strings.Contains(s, "superpassword") {
		t.Errorf("plaintext password found after json marshalling, json is %#v", s)
	}
}

func TestUserRequestNoPassword(t *testing.T) {
	ur := UserRequest{
		User: User{
			Name:     "Mrs. Beakley",
			Email:    "mrsbeakley@yahoo.com",
			Bio:      "Webby Vanderquack's granny",
			Password: "password123",
		},
	}

	b, err := json.Marshal(ur)
	if err != nil {
		t.Error("failed to marshal UserRequest json", err)
	}

	s := string(b)

	if strings.Contains(s, "password123") {
		t.Errorf("plaintext password found after json marshalling, json is %#v", s)
	}
}
