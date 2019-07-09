// +build integration

package user

import (
	"testing"

	"github.com/go-test/deep"

	"github.com/dzeban/conduit/app"
)

//
// XXX: Integration test environment initialized in service_test.go
//

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		email string
		user  *app.User
		err   error
	}{
		{"valid", testUser.Email, testUser, nil},
		{"notfound", "nosuchuser@example.com", nil, app.ErrUserNotFound},
	}

	for _, expected := range tests {
		t.Run(expected.name, func(t *testing.T) {
			u, err := testService.Get(expected.email)

			if err != expected.err {
				t.Errorf("invalid error: expected %v got %v", expected.err, err)
			}

			if diff := deep.Equal(u, expected.user); diff != nil {
				t.Errorf("invalid user: %v", diff)
			}
		})
	}
}

func TestProfile(t *testing.T) {
	profile := &app.Profile{
		Name:  testUser.Name,
		Bio:   testUser.Bio,
		Image: testUser.Image,
	}

	tests := []struct {
		name    string
		profile *app.Profile
		err     error
	}{
		{profile.Name, profile, nil},
		{"nosuchuser", nil, app.ErrUserNotFound},
	}

	for _, expected := range tests {
		p, err := testService.Profile(expected.name)

		if err != expected.err {
			t.Errorf("invalid error: expected %v got %v", expected.err, err)
		}

		if diff := deep.Equal(p, expected.profile); diff != nil {
			t.Errorf("invalid user: %v", diff)
		}
	}
}

func TestLogin(t *testing.T) {
	// New user can't login
	newUser := app.User{
		Name:     "new",
		Email:    "new@example.com",
		Password: "new",
	}

	userValid := *testUser
	userValid.Password = "test"

	userInvalid := *testUser
	userInvalid.Password = "invalidpassword"

	tests := []struct {
		name   string
		user   app.User
		hasErr bool
	}{
		{"valid", userValid, false},
		{"invalid", userInvalid, true},
		{"new", newUser, true},
	}

	for _, expected := range tests {
		t.Run(expected.name, func(t *testing.T) {
			_, err := testService.Login(app.UserRequest{User: expected.user})

			hasErr := (err != nil)
			if hasErr != expected.hasErr {
				t.Errorf("err doesn't match, expected '%v', got '%v', user is %#v\n", expected.hasErr, hasErr, expected.user)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	// Prepare test user
	userToUpdate := app.User{
		Name:     "toupdate",
		Email:    "toupdate@example.com",
		Password: "toupdate",
	}

	u, err := testService.Register(app.UserRequest{User: userToUpdate})
	if err != nil {
		t.Error("failed to prepare test user: " + err.Error())
	}

	// Prepare update info
	u.Name = "newname"
	u.Bio = "mystery"
	u.Image = "image_in_base64"

	// Remove password field to prevent password update
	savedHash := u.Password
	u.Password = ""

	updatedUser, err := testService.Update(u.Email, app.UserRequest{User: *u})
	if err != nil {
		t.Error("failed to get user after update: " + err.Error())
	}

	// Restore password to correctly compare objects
	u.Password = savedHash

	if diff := deep.Equal(u, updatedUser); diff != nil {
		t.Errorf("users not matched after update: %v", diff)
	}

	getUser, err := testService.Get(u.Email)
	if err != nil {
		t.Error("failed to get user after update: " + err.Error())
	}

	if diff := deep.Equal(u, getUser); diff != nil {
		t.Errorf("users not matched after update: %v", diff)
	}
}

// TestUpdatePassword updates password and then try to login with the new one
func TestUpdatePassword(t *testing.T) {
	oldPassword := "old_password"
	newPassword := "new_password"

	// Prepare test user
	userToUpdate := app.User{
		Name:     "passwordupdate",
		Email:    "passwordupdate@example.com",
		Password: oldPassword,
	}

	u, err := testService.Register(app.UserRequest{User: userToUpdate})
	if err != nil {
		t.Error("failed to prepare test user: " + err.Error())
	}

	// Update password
	u.Password = newPassword

	updatedUser, err := testService.Update(u.Email, app.UserRequest{User: *u})
	if err != nil {
		t.Error("failed to get user after update: " + err.Error())
	}

	// Set plain-text password field because update return password hash
	updatedUser.Password = newPassword

	// Try to login with the new password
	_, err = testService.Login(app.UserRequest{User: *updatedUser})
	if err != nil {
		t.Error(err)
	}
}
