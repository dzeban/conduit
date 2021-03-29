package app

// Profile is a public user info with restricted set of fields
type Profile struct {
	Id        int    `json:"id"`
	Name      string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`     // base64 encoded
	Following bool   `json:"following"` // set for authenticated users
}

// ProfileResponse is a structure returned in profile handlers
type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

// ProfileFromUser converts User type to Profile.
// It sets Following field to false because it must be calculated by caller.
func ProfileFromUser(u *User) *Profile {
	if u == nil {
		return nil
	}

	return &Profile{
		Id:        u.Id,
		Name:      u.Name,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: false,
	}
}
