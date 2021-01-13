package user

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/password"
)

type UpdateRequest struct {
	User UpdateUser `json:"user"`
}

type UpdateUser struct {
	Name     string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"` // NOTE: Plaintext password from user
	Bio      string `json:"bio,omitempty"`
	Image    string `json:"image,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	if r.User.Name == "" &&
		r.User.Email == "" &&
		r.User.Bio == "" &&
		r.User.Image == "" &&
		r.User.Password == "" {
		return errors.New("at least one of username, email, bio, image, password is required for update")
	}

	return nil
}

// Update modifies user found by id with the new data passed in user.
// It returns updated user.
func (s *Service) Update(id int, req *UpdateRequest) (*app.User, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ValidationError(err)
	}

	// Check user exists
	u, err := s.store.GetUserById(id)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get user for update"))
	}

	if u == nil {
		return nil, app.ServiceError(app.ErrorUserNotFound)

	}

	// Set fields to update
	u.Name = req.User.Name
	u.Email = req.User.Email
	u.Bio = req.User.Bio
	u.Image = req.User.Image

	// If password is being changed, make the hash from it
	if req.User.Password != "" {
		hash, err := password.HashAndEncode(req.User.Password)
		if err != nil {
			return nil, app.InternalError(errors.Wrap(err, "failed to create password hsah"))
		}

		u.PasswordHash = hash
	}

	err = s.store.UpdateUser(u)
	if err != nil {
		return nil, app.ServiceError(errors.Wrap(err, "failed to update user"))
	}

	// Return updated user
	return s.store.GetUserById(id)
}
