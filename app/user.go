package app

// User represents a user
// User is identified by email and authenticated by JWT
type User struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio,omitempty"`
	Image    []byte `json:"image,omitempty"` // base64 encoded
	Token    string `json:"token,omitempty"`
	Password string `json:"-"`
}

// UserRequest represent a json structure used
// in user serice requests and responses.
type UserRequest struct {
	User User `json:"user"`
}

// UserService defines and interface to work with users
type UserService interface {
	// Login(name, password string) (*User, error)
	Register(u User, password string) error
	// Get(email string) (*User, error)
	// Update(email string, newData User) (*User, error)
}

// UsersConfig describes configuration for UserService
type UsersConfig struct {
	Type string
}
