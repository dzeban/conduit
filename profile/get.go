package profile

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

func (s *Service) Get(username string, currentUser *app.User) (*app.Profile, error) {
	p, err := s.store.GetProfile(username, app.ProfileFromUser(currentUser))
	if err == app.ErrorProfileNotFound {
		return nil, app.ServiceError(app.ErrorProfileNotFound)
	}
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to get profile"))
	}

	return p, nil
}
