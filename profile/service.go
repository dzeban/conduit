package profile

import "github.com/dzeban/conduit/app"

type Service struct {
	store Store
}

type Store interface {
	GetProfile(username string, follower *app.Profile) (*app.Profile, error)
	FollowProfile(follower, followee *app.Profile) error
	UnfollowProfile(follower, followee *app.Profile) error
}

func NewService(store Store) *Service {
	return &Service{store}
}
