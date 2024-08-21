package check

import (
	"testing"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface{}

// Suite registers the given value as a test suite to be run. Any methods
// starting with the Test prefix in the given value will be considered as
// a test method.
func Suite(suite interface{}) interface{} {
	allSuites = append(allSuites, suite)
	return suite
}

// -----------------------------------------------------------------------
// Public running interface.

// TestingT runs all test suites registered with the Suite function,
// printing results to stdout, and reporting any failures back to
// the "testing" package.
func TestingT(t *testing.T) {
	t.Helper()
	RunAll(t)
}

// RunAll runs all test suites registered with the Suite function, using the
// provided run configuration.
func RunAll(t *testing.T) {
	t.Helper()
	for _, suite := range allSuites {
		Run(t, suite)
	}
}

// Run runs the provided test suite using the provided run configuration.
func Run(t *testing.T, suite interface{}) {
	t.Helper()
	runner := newSuiteRunner(suite)
	runner.run(t)
}
