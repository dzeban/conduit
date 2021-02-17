package mock

import "github.com/dzeban/conduit/app"

var (
	Profile1 = app.Profile{
		Id:   1,
		Name: "test",
		Bio:  "test bio",
	}
	Profile2 = app.Profile{
		Id:   2,
		Name: "new",
		Bio:  "new bio",
	}
)

// ProfilesStore is a fake implementation of profiles.Store as Go map
type ProfilesStore struct {
	m map[string]app.Profile
}

func NewProfilesStore() *ProfilesStore {
	ps := &ProfilesStore{
		m: make(map[string]app.Profile),
	}

	for _, profile := range []app.Profile{Profile1, Profile2} {
		ps.m[profile.Name] = profile
	}

	return ps
}

func (ps *ProfilesStore) GetProfile(name string) (*app.Profile, error) {
	p, ok := ps.m[name]
	if !ok {
		return nil, app.ErrorProfileNotFound
	}

	return &p, nil
}
