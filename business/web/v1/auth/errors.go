package auth

import (
	"errors"
	"fmt"
	"runtime"
)

// AuthError is used to pass an error during the request through the
// application with auth specific context.
type AuthError struct {
	msg string
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(format string, args ...interface{}) error {
	return &AuthError{
		msg: fmt.Sprintf(format, args...),
	}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (ae *AuthError) Error() string {
	return ae.msg
}

// IsAuthError checks if an error of type AuthError exists.
func IsAuthError(err error) bool {
	var ae *AuthError
	return errors.As(err, &ae)
}

// WrapError wraps the provided error with file and line number information.
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New("error occurred, but caller info could not be retrieved")
	}
	return fmt.Errorf("%w at %s:%d", err, file, line)
}
