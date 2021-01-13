package user

import (
	"errors"
	"testing"

	"github.com/go-test/deep"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
	"github.com/dzeban/conduit/password"
)

func TestUpdate(t *testing.T) {
	// Create update users to compare in tests
	newUsername := "admin"
	newUserUpdatedUsername := mock.UserUpdatedUsername
	newUserUpdatedUsername.Name = newUsername

	// Test cases
	tests := []struct {
		name    string
		id      int
		req     *UpdateRequest
		errType app.ErrorType
		err     error
		user    *app.User
	}{
		{
			"EmptyValidation",
			mock.UserValid.Id,
			&UpdateRequest{},
			app.ErrorTypeValidation,
			nil,
			nil,
		},
		{
			"AbsentUser",
			-1,
			&UpdateRequest{UpdateUser{Bio: "blah"}},
			app.ErrorTypeService,
			nil,
			nil,
		},
		{
			"UpdateUsername",
			mock.UserUpdatedUsername.Id,
			&UpdateRequest{
				UpdateUser{
					Name: newUsername,
				},
			},
			0,
			nil,
			&newUserUpdatedUsername,
		},
	}

	s := NewService(mock.NewUserStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := s.Update(tt.id, tt.req)
			if err != nil {
				var e app.Error
				// Unwrap service.Error
				if !errors.As(err, &e) {
					t.Errorf("Update(%v): invalid error: expected %T, got %T", tt.req, e, err)
					return
				}

				// Check error type
				if e.Type != tt.errType {
					t.Errorf("Update(%v): invalid error type: expected %v, got %v", tt.req, e.Type, tt.errType)
					return
				}

				// Check error value
				if tt.err != nil {
					if e.Err != tt.err {
						t.Errorf("Update(%v): invalid error value: expected %v, got %v", tt.req, e.Err, tt.err)
						return
					}

				}
			}

			// Check returned user
			if tt.user != nil {
				if diff := deep.Equal(u, tt.user); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	store := mock.NewUserStore()

	// Populate test data
	const oldPassword = "test"
	hash, err := password.HashAndEncode(oldPassword)
	if err != nil {
		t.Fatal(err)
	}

	userUpdatedPassword := app.User{
		Email:        "password_update@example.com",
		Name:         "password_update",
		PasswordHash: hash,
	}

	_ = store.AddUser(&userUpdatedPassword)

	newPassword := "qwerty"

	req := &UpdateRequest{
		UpdateUser{
			Password: newPassword,
		},
	}

	s := NewService(store)
	u, err := s.Update(userUpdatedPassword.Id, req)
	if err != nil {
		t.Errorf("Update(%v): unexpected error: %v", req, err)
	}

	// Check password hash using package function because hash salt is generated
	// randomly
	ok, err := password.Check(newPassword, u.PasswordHash)
	if err != nil {
		t.Errorf("Update(%v): unexpected error on password check: %v", req, err)
	}

	if !ok {
		t.Errorf("Update(%v): password wasn't updated", req)
	}
}
