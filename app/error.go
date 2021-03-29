package app

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrorInternal = errors.New("internal error")
	ErrorLogin    = errors.New("failed to login")
	ErrorRegister = errors.New("failed to register")
	ErrorUpdate   = errors.New("failed to update")

	ErrorUserNotInContext    = errors.New("no user in context")
	ErrorUserNotFound        = errors.New("user not found")
	ErrorUserExists          = errors.New("user exists")
	ErrorUserUpdateForbidden = errors.New("user update forbidden")
	ErrorPasswordMismatch    = errors.New("password mismatch")

	ErrorProfileNotFound            = errors.New("profile not found")
	ErrorProfileAlreadyFollowing    = errors.New("profile is already followed")
	ErrorProfileAlreadyNotFollowing = errors.New("profile is already not followed")

	ErrorArticleExists          = errors.New("article exists")
	ErrorArticleNotFound        = errors.New("article not found")
	ErrorArticleUpdateForbidden = errors.New("article update forbidden")
	ErrorArticleDeleteForbidden = errors.New("article delete forbidden")
	ErrorArticleInvalidLimit    = errors.New("invalid article list limit")
	ErrorArticleInvalidOffset   = errors.New("invalid article list offset")
	ErrorArticleInvalidFilter   = errors.New("invalid article list filter")

	ErrorValidationUsernameIsRequired = errors.New("username is required")
	ErrorValidationEmailIsRequired    = errors.New("email is required")
	ErrorValidationPasswordIsRequired = errors.New("password is required")
	ErrorValidationTitleIsRequired    = errors.New("title is required")
	ErrorValidationBodyIsRequired     = errors.New("body is required")
)

type ErrorType int

const (
	ErrorTypeInternal = iota
	ErrorTypeService
	ErrorTypeValidation
)

func (et ErrorType) String() string {
	switch et {
	case ErrorTypeInternal:
		return "ErrorTypeInternal"
	case ErrorTypeService:
		return "ErrorTypeService"
	case ErrorTypeValidation:
		return "ErrorTypeValidation"
	default:
		return fmt.Sprintf("%d", et)
	}
}

type Error struct {
	Type ErrorType
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) Unwrap() error {
	return e.Err
}

func InternalError(err error) Error {
	return Error{ErrorTypeInternal, err}
}

func ServiceError(err error) Error {
	return Error{ErrorTypeService, err}
}

func ValidationError(err error) Error {
	return Error{ErrorTypeValidation, err}
}

func ServerError(err error, msg string) string {
	// errors.Wrap doesn't handle nil errors. To avoid nil pointer in error
	// message we create empty error here when error is nil
	if err == nil {
		err = errors.New("")
	}
	return fmt.Sprintf(`{"error":{"message":["%s"]}}`, errors.Wrap(err, msg))
}
