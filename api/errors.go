package api

import "fmt"

type ErrorType int

const (
	_ ErrorType = iota
	errInvalidID
	errUnauthorized
	errInternal
	errInvalidInput
	errNotFound
	errAlreadyExists
	errForbidden
	errBadRequest
)

type AppError struct {
	ErrType ErrorType
	Code    int
	Err     error
}

var (
	InvalidID     = AppError{ErrType: errInvalidID}
	Unauthorized  = AppError{ErrType: errUnauthorized}
	Internal      = AppError{ErrType: errInternal}
	InvalidInput  = AppError{ErrType: errInvalidInput}
	NotFound      = AppError{ErrType: errNotFound}
	AlreadyExists = AppError{ErrType: errAlreadyExists}
	Forbidden     = AppError{ErrType: errForbidden}
	BadRequest    = AppError{ErrType: errBadRequest}
)

func (appErr AppError) Error() string {
	switch appErr.ErrType {
	case errInvalidID:
		return "The provided ID is invalid"
	case errUnauthorized:
		return fmt.Sprintf("Access denied: %d Unauthorized", appErr.Code)
	case errInternal:
		return "Internal server error"
	case errInvalidInput:
		return "Invalid input"
	case errNotFound:
		return fmt.Sprintf("%d Resource not found", appErr.Code)
	case errAlreadyExists:
		return fmt.Sprintf("%d Resource already exists", appErr.Code)
	case errForbidden:
		return fmt.Sprintf("%d Forbidden", appErr.Code)
	case errBadRequest:
		return fmt.Sprintf("%d Bad request", appErr.Code)
	default:
		return "Unknown error"
	}
}

//===============
// ERROR HELPERS
//===============

// with returns an error with a particular code (e.g unauthorized)
func (appErr AppError) with(code int) AppError {
	err := appErr
	err.Code = code
	return err
}

func (appErr AppError) from(location string, err error) AppError {
	errFrom := appErr
	errFrom.Err = fmt.Errorf("from %s: %v", location, err)
	return errFrom
}

func (appErr *AppError) Unwrap() error {
	return appErr.Err
}

func (appErr *AppError) Is(target error) bool {
	err, ok := target.(*AppError) // reflection to check if target is an AppError

	if !ok {
		return false
	}

	return appErr.ErrType == err.ErrType
}
