package user

import (
	"errors"
	"testing"

	"github.com/dzeban/conduit/app"
	"github.com/dzeban/conduit/mock"
)

func TestRegister(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		req     *RegisterRequest
		errType app.ErrorType
		err     error
	}{
		{
			"EmptyValidation",
			&RegisterRequest{},
			app.ErrorTypeValidation,
			nil,
		},
		{
			"EmailValidation",
			&RegisterRequest{
				RegisterUser{
					Username: mock.UserValid.Name,
					Password: mock.TestPassword,
				},
			},
			app.ErrorTypeValidation,
			app.ErrorValidationEmailIsRequired,
		},
		{
			"PasswordValidation",
			&RegisterRequest{
				RegisterUser{
					Email:    mock.UserValid.Email,
					Username: mock.UserValid.Name,
				},
			},
			app.ErrorTypeValidation,
			app.ErrorValidationPasswordIsRequired,
		},
		{
			"UsernameValidation",
			&RegisterRequest{
				RegisterUser{
					Email:    mock.UserValid.Email,
					Password: mock.TestPassword,
				},
			},
			app.ErrorTypeValidation,
			app.ErrorValidationUsernameIsRequired,
		},
		{
			"UserExists",
			&RegisterRequest{
				RegisterUser{
					Email:    mock.UserValid.Email,
					Username: mock.UserValid.Name,
					Password: mock.TestPassword,
				},
			},
			app.ErrorTypeService,
			app.ErrorUserExists,
		},
		{
			"Valid",
			&RegisterRequest{
				RegisterUser{
					Email:    "new@example.com",
					Username: "new",
					Password: "new",
				},
			},
			0,
			nil,
		},
	}

	s := NewService(mock.NewUserStore())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Register(tt.req)

			// Check error
			if err != nil {
				var e app.Error
				// Unwrap service.Error
				if !errors.As(err, &e) {
					t.Errorf("Register(%+v): invalid error: expected %T, got %T", tt.req, e, err)
					return
				}

				// Check error type
				if e.Type != tt.errType {
					t.Errorf("Register(%+v): invalid error type: expected %v, got %v", tt.req, tt.errType, e.Type)
					return
				}

				// Check error value
				if tt.err != nil {
					if e.Err != tt.err {
						t.Errorf("Register(%+v): invalid error value: expected %v, got %v", tt.req, tt.err, e.Err)
						return
					}

				}
			}
		})
	}
}
