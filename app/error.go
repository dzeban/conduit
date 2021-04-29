package app

import (
	"fmt"
)

type ErrorType int

const (
	ErrorTypeInternal = iota
	ErrorTypeService
	ErrorTypeAuth
)

func (et ErrorType) String() string {
	switch et {
	case ErrorTypeInternal:
		return "ErrorTypeInternal"
	case ErrorTypeService:
		return "ErrorTypeService"
	case ErrorTypeAuth:
		return "ErrorTypeAuth"
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

func AuthError(err error) Error {
	return Error{ErrorTypeAuth, err}
}
