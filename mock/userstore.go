package mock

import (
	"github.com/pkg/errors"

	"github.com/dzeban/conduit/app"
)

const (
	TestPassword = "test"

	// Encoded hash for "test". Created with `./hashpass test`
	TestPasswordHash = "$argon2id$v=19$m=32768,t=5,p=1$5rGEEAspznZZGSRgww27Bg$iBaYNvlqwdYgjOLr7rXupS+/YUuAMemUzazttfu9POQsZymmTiUrRAqxjyQA6dKGUMsuAqbu5IukmZaJNmUV8w"
)

var (
	UserValid = app.User{
		Email:        "test@example.com",
		Name:         "test",
		PasswordHash: TestPasswordHash,
	}

	UserUpdatedUsername = app.User{
		Email:        "updated@example.com",
		Name:         "updated",
		PasswordHash: TestPasswordHash,
	}

	UserInvalid = app.User{
		Email:        "invalid_hash@example.com",
		Name:         "invalid_hash",
		PasswordHash: "xxx",
	}
)

// UserStore is a fake implementation of user.Store as Go map
type UserStore struct {
	m map[string]*app.User
}

func NewUserStore() *UserStore {
	us := &UserStore{
		m: make(map[string]*app.User),
	}

	_ = us.AddUser(&UserValid)
	_ = us.AddUser(&UserUpdatedUsername)
	_ = us.AddUser(&UserInvalid)

	return us
}

func (us *UserStore) GetUser(email string) (*app.User, error) {
	return us.m[email], nil
}

func (us *UserStore) AddUser(user *app.User) error {
	us.m[user.Email] = user
	return nil
}

func (us *UserStore) UpdateUser(email string, user *app.User) error {
	u := us.m[email]

	if u == nil {
		return errors.New("user not found")
	}

	if user.Name != "" {
		u.Name = user.Name
	}

	if user.Bio != "" {
		u.Bio = user.Bio
	}

	if user.Image != "" {
		u.Image = user.Image
	}

	if user.PasswordHash != "" {
		u.PasswordHash = user.PasswordHash
	}

	// If we update email, recreate user under the new key
	if user.Email != "" && email != user.Email {
		delete(us.m, email)
		u.Email = user.Email
	}

	us.m[u.Email] = u

	return nil
}
