package web

import (
	"errors"
	"fmt"
	"runtime"
)

// wrapError wraps the provided error with file and line number information.
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
