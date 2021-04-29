package profile

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

func (s *Service) Follow(follower *app.User, username string) (*app.Profile, error) {
	followee, err := s.Get(username, follower)
	if err != nil {
		return nil, app.ServiceError(errors.Wrap(err, "failed to get followee profile"))
	}

	if followee.Following {
		return nil, app.ServiceError(errorProfileAlreadyFollowing)
	}

	err = s.store.FollowProfile(app.ProfileFromUser(follower), followee)
	if err != nil {
		return nil, app.InternalError(errors.Wrap(err, "failed to follow profile"))
	}

	p, err := s.Get(username, follower)
	if err != nil {
		return nil, app.InternalError(app.ErrorProfileNotFound)
	}

	return p, nil
}
