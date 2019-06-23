package user

import (
	"encoding/json"

	"github.com/dzeban/conduit/app"
)

// Redeclare user from app.User to allow password serialization
// app.User type has special MarshalJSON that disables password serialization to
// prevent accidental password exposure.
// This type overrides this behavior by creating alias and implementing standard
// json marshalling.
type user app.User

func (u user) MarshalJSON() ([]byte, error) {
	// Use custom aliased type here to avoid
	// infinite recursion in marshalling
	type alias user
	return json.Marshal(struct {
		alias
	}{
		alias: (alias)(u),
	})
}

type userRequest struct {
	User user `json:"user"`
}
