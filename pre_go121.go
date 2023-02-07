//go:build !go1.21
// +build !go1.21

package check

// nilIfPanicNilError returns nil if v is an instance of runtime.PanicNilError (introduced in go1.21)
func nilIfPanicNilError(v interface{}) interface{} {
	return v
}
