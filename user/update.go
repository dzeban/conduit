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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // NOTE: Plaintext password from user
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func (r *UpdateRequest) Validate() error {
	if r.User.Username == "" &&
		r.User.Email == "" &&
		r.User.Bio == "" &&
		r.User.Image == "" &&
		r.User.Password == "" {
		return errors.New("at least one of username, email, bio, image, password is required for update")
	}

	return nil
}

// Update modifies user found by email with the new data passed in user.
// It returns updated user.
func (s *Service) Update(email string, req *UpdateRequest) (*app.User, error) {
	// Validate request
	err := req.Validate()
	if err != nil {
		return nil, app.ValidationError(err)
	}

	// Set fields to update
	user := &app.User{
		Name:  req.User.Username,
		Email: req.User.Email,
		Bio:   req.User.Bio,
		Image: req.User.Image,
	}

	// If password is being changed, make the hash from it
	if req.User.Password != "" {
		hash, err := password.HashAndEncode(req.User.Password)
		if err != nil {
			return nil, app.InternalError(errors.Wrap(err, "failed to create password hsah"))
		}

		user.PasswordHash = hash
	}

	err = s.store.UpdateUser(email, user)
	if err != nil {
		return nil, app.ServiceError(errors.Wrap(err, "failed to update user"))
	}

	// Return updated user
	return s.store.GetUser(email)
}
