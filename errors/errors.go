//IDEA from https://hackernoon.com/golang-handling-errors-gracefully-8e27f1db729f

package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorType string

const (
	NoType = ErrorType(iota)
	BadRequest
	NotFound
	Unauthorized
	Forbidden
	PreconditionFailed
	InternalServerError
	ServiceUnavailable
)

type CustomError struct {
	errorType     ErrorType
	originalError error
}

func (error CustomError) Error() string {
	return error.originalError.Error()
}

func (errType ErrorType) New(msg string) error {
	return CustomError{errType, errors.New(msg)}
}

func (errType ErrorType) Newf(msg string, args ...interface{}) error {
	return CustomError{errType, fmt.Errorf(msg, args...)}
}

func (errType ErrorType) Wrap(err error, msg string) error {
	return errType.Wrapf(err, msg)
}

func (errType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	wrappedErr := errors.Wrapf(err, msg, args)
	return CustomError{errType, wrappedErr}
}

func New(msg string) error {
	return CustomError{NoType, errors.New(msg)}
}

func Newf(msg string, args ...interface{}) error {
	return CustomError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args)
	if customErr, ok := err.(CustomError); ok {
		return CustomError{customErr.errorType, wrappedError}
	}
	return CustomError{NoType, wrappedError}
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(CustomError); ok {
		return customErr.errorType
	}

	return NoType
}
