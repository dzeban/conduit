package user

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
)

func TestLogin(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		req     *LoginRequest
		errType app.ErrorType
		err     error
	}{
		{
			"EmptyValidation",
			&LoginRequest{},
			app.ErrorTypeValidation,
			nil,
		},
		{
			"EmailValidation",
			&LoginRequest{
				LoginUser{
					Password: mock.TestPassword,
				},
			},
			app.ErrorTypeValidation,
			app.ErrorValidationEmailIsRequired,
		},
		{
			"PasswordValidation",
			&LoginRequest{
				LoginUser{
					Email: mock.UserValid.Email,
				},
			},
			app.ErrorTypeValidation,
			app.ErrorValidationPasswordIsRequired,
		},
		{
			"NonExist",
			&LoginRequest{
				LoginUser{
					Email:    "no_such_user@example.com",
					Password: "abc",
				},
			},
			app.ErrorTypeService,
			app.ErrorUserNotFound,
		},
		{
			"InvalidPassword",
			&LoginRequest{
				LoginUser{
					Email:    mock.UserValid.Email,
					Password: "invalid",
				},
			},
			app.ErrorTypeService,
			app.ErrorPasswordMismatch,
		},
		{
			"InvalidPasswordHash",
			&LoginRequest{
				LoginUser{
					Email:    mock.UserInvalid.Email,
					Password: "some",
				},
			},
			app.ErrorTypeInternal,
			nil,
		},
		{
			"Valid",
			&LoginRequest{
				LoginUser{
					Email:    mock.UserValid.Email,
					Password: mock.TestPassword,
				},
			},
			0,
			nil,
		},
	}

	s := NewService(mock.NewUserStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Login(tt.req)
			if err != nil {
				var e app.Error
				// Unwrap service.Error
				if !errors.As(err, &e) {
					t.Errorf("Login(%v): invalid error: expected %T, got %T", tt.req, e, err)
					return
				}

				// Check error type
				if e.Type != tt.errType {
					t.Errorf("Login(%v): invalid error type: expected %v, got %v", tt.req, tt.errType, e.Type)
					return
				}

				// Check error value
				if tt.err != nil {
					if e.Err != tt.err {
						t.Errorf("Login(%v): invalid error value: expected %v, got %v", tt.req, tt.err, e.Err)
						return
					}

				}
			}
		})
	}
}
