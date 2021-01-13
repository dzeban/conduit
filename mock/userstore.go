package mock

import (
	"math/rand"

	"github.com/dzeban/conduit/app"
)

const (
	TestPassword = "test"

	// Encoded hash for "test". Created with `./hashpass test`
	TestPasswordHash = "$argon2id$v=19$m=32768,t=5,p=1$5rGEEAspznZZGSRgww27Bg$iBaYNvlqwdYgjOLr7rXupS+/YUuAMemUzazttfu9POQsZymmTiUrRAqxjyQA6dKGUMsuAqbu5IukmZaJNmUV8w"
)

var (
	UserValid = app.User{
		Id:           1,
		Email:        "test@example.com",
		Name:         "test",
		PasswordHash: TestPasswordHash,
	}

	UserUpdatedUsername = app.User{
		Id:           2,
		Email:        "updated@example.com",
		Name:         "updated",
		PasswordHash: TestPasswordHash,
	}

	UserInvalid = app.User{
		Id:           3,
		Email:        "invalid_hash@example.com",
		Name:         "invalid_hash",
		PasswordHash: "xxx",
	}
)

// UserStore is a fake implementation of user.Store as Go map
type UserStore struct {
	ById    map[int]app.User
	ByEmail map[string]app.User
}

func NewUserStore() *UserStore {
	us := &UserStore{
		ById:    make(map[int]app.User),
		ByEmail: make(map[string]app.User),
	}

	for _, user := range []*app.User{&UserValid, &UserUpdatedUsername, &UserInvalid} {
		err := us.AddUser(user)
		if err != nil {
			panic(err)
		}
	}

	return us
}

func (us *UserStore) GetUser(email string) (*app.User, error) {
	u, ok := us.ByEmail[email]
	if !ok {
		return nil, nil
	}
	return &u, nil
}

func (us *UserStore) GetUserById(id int) (*app.User, error) {
	u, ok := us.ById[id]
	if !ok {
		return nil, nil
	}
	return &u, nil
}

func (us *UserStore) AddUser(user *app.User) error {
	if user.Id == 0 {
		user.Id = 10 + rand.Int() // "10 + " is needed to avoid overlap with predefined mock users
	}
	us.ById[user.Id] = *user
	us.ByEmail[user.Email] = *user

	return nil
}

func (us *UserStore) UpdateUser(newUser *app.User) error {
	user, ok := us.ById[newUser.Id]
	if !ok {
		return app.ErrorUserNotFound
	}

	if newUser.Name != "" {
		user.Name = newUser.Name
	}

	if newUser.Bio != "" {
		user.Bio = newUser.Bio
	}

	if newUser.Image != "" {
		user.Image = newUser.Image
	}

	if newUser.PasswordHash != "" {
		user.PasswordHash = newUser.PasswordHash
	}

	// If we update email, recreate user under the new key
	if newUser.Email != "" && user.Email != newUser.Email {
		delete(us.ByEmail, user.Email)
		user.Email = newUser.Email
	}

	us.ByEmail[user.Email] = user
	us.ById[user.Id] = user

	return nil
}
