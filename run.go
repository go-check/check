package check

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

var (
	oldListFlag = flag.Bool("gocheck.list", false, "List the names of all tests that will be run")
	oldWorkFlag = flag.Bool("gocheck.work", false, "Display and do not remove the test working directory")

	newListFlag = flag.Bool("check.list", false, "List the names of all tests that will be run")
	newWorkFlag = flag.Bool("check.work", false, "Display and do not remove the test working directory")
)

// TestingT runs all test suites registered with the Suite function,
// printing results to stdout, and reporting any failures back to
// the "testing" package.
func TestingT(t *testing.T) {
	t.Helper()
	conf := &RunConf{
		KeepWorkDir: *oldWorkFlag || *newWorkFlag,
	}
	if *oldListFlag || *newListFlag {
		w := bufio.NewWriter(os.Stdout)
		for _, name := range ListAll(conf) {
			fmt.Fprintln(w, name)
		}
		w.Flush()
		return
	}
	RunAll(t, conf)
}

// RunAll runs all test suites registered with the Suite function, using the
// provided run configuration.
func RunAll(t *testing.T, runConf *RunConf) {
	t.Helper()
	for _, suite := range allSuites {
		t.Run(suiteName(suite), func(t *testing.T) {
			Run(t, suite, runConf)
		})
	}
}

// Run runs the provided test suite using the provided run configuration.
func Run(t *testing.T, suite interface{}, runConf *RunConf) {
	t.Helper()
	runner := newSuiteRunner(suite, runConf)
	runner.run(t)
}

// ListAll returns the names of all the test functions registered with the
// Suite function that will be run with the provided run configuration.
func ListAll(runConf *RunConf) []string {
	var names []string
	for _, suite := range allSuites {
		names = append(names, List(suite, runConf)...)
	}
	return names
}

// List returns the names of the test functions in the given
// suite that will be run with the provided run configuration.
func List(suite interface{}, runConf *RunConf) []string {
	var names []string
	runner := newSuiteRunner(suite, runConf)
	for _, t := range runner.tests {
		names = append(names, t.String())
	}
	return names
}
