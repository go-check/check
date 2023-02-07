//go:build go1.21
// +build go1.21

package check

import (
	"errors"
	"runtime"
)

// nilIfPanicNilError returns nil if v is an instance of runtime.PanicNilError (introduced in go1.21)
func nilIfPanicNilError(v interface{}) interface{} {
	if err, isError := v.(error); isError {
		if panicNilError := (&runtime.PanicNilError{}); errors.As(err, &panicNilError) {
			return nil
		}
	}
	return v
}
