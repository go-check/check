package gocheck


import (
    "strings"
    "fmt"
)


// -----------------------------------------------------------------------
// Basic succeeding/failing logic.

// Return true if the currently running test has already failed.
func (c *C) Failed() bool {
    return c.status == failedSt
}

// Mark the currently running test as failed. Something ought to have been
// previously logged so that the developer knows what went wrong. The higher
// level helper functions will fail the test and do the logging properly.
func (c *C) Fail() {
    c.status = failedSt
}

// Mark the currently running test as failed, and stop running the test.
// Something ought to have been previously logged so that the developer
// knows what went wrong. The higher level helper functions will fail the
// test and do the logging properly.
func (c *C) FailNow() {
    c.Fail()
    c.stopNow()
}

// Mark the currently running test as succeeded, undoing any previous
// failures.
func (c *C) Succeed() {
    c.status = succeededSt
}

// Mark the currently running test as succeeded, undoing any previous
// failures, and stop running the test.
func (c *C) SucceedNow() {
    c.Succeed()
    c.stopNow()
}

// Expect the currently running test to fail, for the given reason.  If the
// test does not fail, an error will be reported to raise the attention to
// this fact. The reason string is just a summary of why the given test is
// supposed to fail.  This method is useful to temporarily disable tests
// which cover well known problems until a better time to fix the problem
// is found, without forgetting about the fact that a failure still exists.
func (c *C) ExpectFailure(reason string) {
    c.expectedFailure = &reason
}


// -----------------------------------------------------------------------
// Basic logging.

// Return the current test error output.
func (c *C) GetTestLog() string {
    return c.logb.String()
}

// Log some information into the test error output.  The provided arguments
// will be assembled together into a string using fmt.Sprint().
func (c *C) Log(args ...interface{}) {
    c.log(args...)
}

// Log some information into the test error output.  The provided arguments
// will be assembled together into a string using fmt.Sprintf().
func (c *C) Logf(format string, args ...interface{}) {
    c.logf(format, args...)
}

// Log an error into the test error output, and mark the test as failed.
// The provided arguments will be assembled together into a string using
// fmt.Sprint().
func (c *C) Error(args ...interface{}) {
    c.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    c.Fail()
}

// Log an error into the test error output, and mark the test as failed.
// The provided arguments will be assembled together into a string using
// fmt.Sprintf().
func (c *C) Errorf(format string, args ...interface{}) {
    c.logCaller(1, fmt.Sprintf("Error: " + format, args...))
    c.Fail()
}

// Log an error into the test error output, mark the test as failed, and
// stop the test execution. The provided arguments will be assembled
// together into a string using fmt.Sprint().
func (c *C) Fatal(args ...interface{}) {
    c.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    c.FailNow()
}

// Log an error into the test error output, mark the test as failed, and
// stop the test execution. The provided arguments will be assembled
// together into a string using fmt.Sprintf().
func (c *C) Fatalf(format string, args ...interface{}) {
    c.logCaller(1, fmt.Sprint("Error: ", fmt.Sprintf(format, args...)))
    c.FailNow()
}


// -----------------------------------------------------------------------
// Generic checks and assertions based on checkers.

// Verify if the first value matches with the expected value.  What
// matching means is defined by the provided checker. In case they do not
// match, an error will be logged, the test will be marked as failed, and
// the test execution will continue.  Some checkers may not need the expected
// argument (e.g. IsNil).  In either case, any extra arguments provided to
// the function will be logged next to the reported problem when the
// matching fails.  This is a handy way to provide problem-specific hints.
func (c *C) Check(obtained interface{}, checker Checker,
                  args ...interface{}) bool {
    return c.internalCheck("Check", obtained, checker, args...)
}

// Ensure that the first value matches with the expected value.  What
// matching means is defined by the provided checker. In case they do not
// match, an error will be logged, the test will be marked as failed, and
// the test execution will stop.  Some checkers may not need the expected
// argument (e.g. IsNil).  In either case, any extra arguments provided to
// the function will be logged next to the reported problem when the
// matching fails.  This is a handy way to provide problem-specific hints.
func (c *C) Assert(obtained interface{}, checker Checker,
                   args ...interface{}) {
    if !c.internalCheck("Assert", obtained, checker, args...) {
        c.stopNow()
    }
}

func (c *C) internalCheck(funcName string,
                          obtained interface{}, checker Checker,
                          args ...interface{}) bool {
    if checker == nil {
        c.logCaller(2, fmt.Sprintf("%s(obtained, nil!?, ...):", funcName))
        c.logString("Oops.. you've provided a nil checker!")
        goto fail
    }

    // If the last argument is a bug info, extract it out.
    var bug BugInfo
    if len(args) > 0 {
        if gotBug, hasBug := args[len(args)-1].(BugInfo); hasBug {
            bug = gotBug
            args = args[:len(args)-1]
        }
    }

    // Ensure we got the needed number of arguments in expected.  Note that
    // this logic is a bit more complex than it ought to be, mainly because
    // it's leaving the door open to multiple expected values.
    var expectedWanted int
    var expected interface{}
    if checker.NeedsExpectedValue() {
        expectedWanted = 1
    }
    if len(args) == expectedWanted {
        if expectedWanted > 0 {
            expected = args[0]
        }
    } else {
        obtainedName, expectedName := checker.VarNames()
        c.logCaller(2, fmt.Sprintf("%s(%s, %s, >%s<):", funcName, obtainedName,
                                   checker.Name(), expectedName))
        c.logString(fmt.Sprintf("Wrong number of %s args for %s: " +
                                "want %d, got %d", expectedName, checker.Name(),
                                expectedWanted, len(args)))
        goto fail
    }

    // Do the actual check.
    result, error := checker.Check(obtained, expected)
    if !result || error != "" {
        obtainedName, expectedName := checker.VarNames()
        var summary string
        if expectedWanted > 0 {
            summary = fmt.Sprintf("%s(%s, %s, %s):", funcName, obtainedName,
                                  checker.Name(), expectedName)
        } else {
            summary = fmt.Sprintf("%s(%s, %s):", funcName, obtainedName,
                                  checker.Name())
        }
        c.logCaller(2, summary)
        c.logValue(strings.Title(obtainedName), obtained)
        if expectedWanted > 0 {
            c.logValue(strings.Title(expectedName), expected)
        }
        if error != "" {
            c.logString(error)
        } else if bug != nil {
            c.logString(bug.GetBugInfo())
        }
        goto fail
    }
    return true

fail:
    c.logNewLine()
    c.Fail()
    return false
}
