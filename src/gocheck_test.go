// This file contains just a few generic helpers which are used by the
// other test files.
//
// Tests are distributed the following way:
//
//  gochecker_test.go: This file.  Just generic helpers.
//  bootstrap_test.go: Tests breaking the chicken and egg problem of
//                     testing a testing framework with itself.
//    helpers_test.go: Tests for helper methods in *gocheck.T.

package gocheck_test


import (
    "gocheck"
    "runtime"
    "os"
)


// -----------------------------------------------------------------------
// Helper functions.

// Return the file line where it's called.
func getMyLine() int {
    if _, _, line, ok := runtime.Caller(1); ok {
        return line
    }
    return -1
}


// -----------------------------------------------------------------------
// Helper type implementing a basic io.Writer for testing output.

// Type implementing the io.Writer interface for analyzing output.
type String struct {
    value string
}

// The only function required by the io.Writer interface.  Will append
// written data to the String.value string.
func (s *String) Write(p []byte) (n int, err os.Error) {
    s.value += string(p)
    return len(p), nil
}

// Trivial wrapper to test errors happening on a different file
// than the test itself.
func checkEqualWrapper(t *gocheck.T,
                       expected interface{},
                       obtained interface{}) (result bool, line int) {
    return t.CheckEqual(expected, obtained), getMyLine()
}


// -----------------------------------------------------------------------
// Helper suite for testing basic fail behavior.

type FailHelper struct{
    testLine int
}

func (s *FailHelper) TestLogAndFail(t *gocheck.T) {
    s.testLine = getMyLine()-1
    t.Log("Expected failure!")
    t.Fail()
}


// -----------------------------------------------------------------------
// Helper suite for testing basic success behavior.

type SuccessHelper struct{}

func (s *SuccessHelper) TestLogAndSucceed(t *gocheck.T) {
    t.Log("Expected success!")
}
