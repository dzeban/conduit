package app

// User represents a user
// User is identified by email and authenticated by JWT
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Bio      string `json:"bio,omitempty"`
	Image    []byte `json:"image,omitempty"` // base64 encoded
	Token    string `json:"token,omitempty"`
}

// UserService defines and interface to work with users
type UserService interface {
	Login(name, password string) (*User, error)
	Register(u User) (*User, error)
	Get(email string) (*User, error)
	Update(email string, newData User) (*User, error)
}
