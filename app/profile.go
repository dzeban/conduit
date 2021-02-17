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
