package gocheck

import (
    "testing"
    "flag"
    "fmt"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface{}

// Register the given value as a test suite to be run.  Any methods starting
// with the Test prefix in the given value will be considered as a test to
// be run.
func Suite(suite interface{}) interface{} {
    return Suites(suite)[0]
}

// Register all of the provided values as test suites to be run.  Any methods
// starting with the Test prefix in the given values will be considered as a
// test to be run.
func Suites(suites ...interface{}) []interface{} {
    lenAllSuites := len(allSuites)
    lenSuites := len(suites)
    if lenAllSuites + lenSuites > cap(allSuites) {
        newAllSuites := make([]interface{}, (lenAllSuites+lenSuites)*2)
        copy(newAllSuites, allSuites)
        allSuites = newAllSuites
    }
    allSuites = allSuites[0:lenAllSuites+lenSuites]
    for i, suite := range suites {
        allSuites[lenAllSuites+i] = suite
    }
    return suites
}


// -----------------------------------------------------------------------
// Public running interface.

var filterFlag = flag.String("f", "",
                             "Regular expression to select " +
                             "what to run (gocheck)")

// Run all test suites registered with the Suite() function, printing
// results to stdout, and reporting any failures back to the 'testing'
// module.
func TestingT(testingT *testing.T) {
    result := RunAll(&RunConf{Filter: *filterFlag})
    println(result.String())
    if usedDeprecatedChecks {
        println("WARNING: The Check/AssertEqual() and Check/AssertMatch() " +
                "family of functions\n         is deprecated, and will be " +
                "removed in the future.  Assert() and\n         Check() " +
                "are both more comfortable and more powerful.  See the\n" +
                "         documentation at http://labix.org/gocheck for " +
                "more details.\n")
    }
    if !result.Passed() {
        testingT.Fail()
    }
}

// Run all test suites registered with the Suite() function, using the
// given run configuration.
func RunAll(runConf *RunConf) *Result {
    result := Result{}
    for _, suite := range allSuites {
        result.Add(Run(suite, runConf))
    }
    return &result
}

// Run the given test suite using the provided run configuration.
func Run(suite interface{}, runConf *RunConf) *Result {
    runner := newSuiteRunner(suite, runConf)
    return runner.run()
}


// -----------------------------------------------------------------------
// Result methods.

func (r *Result) Add(other *Result) {
    r.Succeeded += other.Succeeded
    r.Skipped += other.Skipped
    r.Failed += other.Failed
    r.Panicked += other.Panicked
    r.FixturePanicked += other.FixturePanicked
    r.Missed += other.Missed
}

func (r *Result) Passed() bool {
    return (r.Failed == 0 && r.Panicked == 0 &&
            r.FixturePanicked == 0 && r.Missed == 0 &&
            r.RunError == nil)
}

func (r *Result) String() string {
    if r.RunError != nil {
        return "ERROR: " + r.RunError.String()
    }

    var value string
    if r.Failed == 0 && r.Panicked == 0 && r.FixturePanicked == 0 &&
       r.Missed == 0 {
        value = "OK: "
    } else {
        value = "OOPS: "
    }
    value += fmt.Sprintf("%d passed", r.Succeeded)
    if r.Skipped != 0 {
        value += fmt.Sprintf(", %d skipped", r.Skipped)
    }
    if r.Failed != 0 {
        value += fmt.Sprintf(", %d FAILED", r.Failed)
    }
    if r.Panicked != 0 {
        value += fmt.Sprintf(", %d PANICKED", r.Panicked)
    }
    if r.FixturePanicked != 0 {
        value += fmt.Sprintf(", %d FIXTURE PANICKED", r.FixturePanicked)
    }
    if r.Missed != 0 {
        value += fmt.Sprintf(", %d MISSED", r.Missed)
    }
    return value
}
