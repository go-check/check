package gocheck


import (
    "reflect"
    "regexp"
    "fmt"
    "os"
)


// Type passed as an argument to test methods.
type T struct {
    *call
}


// -----------------------------------------------------------------------
// Basic succeeding/failing logic.

// Return true if the currently running test has already failed.
func (t *T) Failed() bool {
    return t.call.status == failedSt
}

// Mark the currently running test as failed. Something ought to have been
// previously logged so that the developer knows what went wrong. The higher
// level helper functions will fail the test and do the logging properly.
func (t *T) Fail() {
    t.call.status = failedSt
}

// Mark the currently running test as failed, and stop running the test.
// Something ought to have been previously logged so that the developer
// knows what went wrong. The higher level helper functions will fail the
// test and do the logging properly.
func (t *T) FailNow() {
    t.Fail()
    t.stopNow()
}

// Mark the currently running test as succeeded, undoing any previous
// failures.
func (t *T) Succeed() {
    t.call.status = succeededSt
}

// Mark the currently running test as succeeded, undoing any previous
// failures, and stop running the test.
func (t *T) SucceedNow() {
    t.Succeed()
    t.stopNow()
}

// Expect the currently running test to fail, for the given reason.  If the
// test does not fail, an error will be reported to raise the attention to
// this fact. The reason string is just a summary of why the given test is
// supposed to fail.  This method is useful to temporarily disable tests
// which cover well known problems until a better time to fix the problem
// is found, without forgetting about the fact that a failure still exists.
func (t *T) ExpectFailure(reason string) {
    t.expectedFailure = &reason
}


// -----------------------------------------------------------------------
// Basic logging.

// Return the current test error output.
func (t *T) GetTestLog() string {
    return t.call.logv
}

// Log some information into the test error output.  The provided arguments
// will be assembled together into a string using fmt.Sprint().
func (t *T) Log(args ...interface{}) {
    t.log(args...)
}

// Log some information into the test error output.  The provided arguments
// will be assembled together into a string using fmt.Sprintf().
func (t *T) Logf(format string, args ...interface{}) {
    t.logf(format, args...)
}

// Log an error into the test error output, and mark the test as failed.
// The provided arguments will be assembled together into a string using
// fmt.Sprint().
func (t *T) Error(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    t.Fail()
}

// Log an error into the test error output, and mark the test as failed.
// The provided arguments will be assembled together into a string using
// fmt.Sprintf().
func (t *T) Errorf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprintf("Error: " + format, args...))
    t.Fail()
}

// Log an error into the test error output, mark the test as failed, and
// stop the test execution. The provided arguments will be assembled
// together into a string using fmt.Sprint().
func(t *T) Fatal(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    t.FailNow()
}

// Log an error into the test error output, mark the test as failed, and
// stop the test execution. The provided arguments will be assembled
// together into a string using fmt.Sprintf().
func(t *T) Fatalf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprintf(format, args...)))
    t.FailNow()
}


// -----------------------------------------------------------------------
// Equality testing.

// Verify if the first value is equal to the second value.  In case
// they're not equal, an error will be logged, the test will be marked as
// failed, and the test execution will continue.  The extra arguments are
// optional and, if provided, will be assembled together with fmt.Sprint()
// and printed next to the reported problem in case of errors. The returned
// value will be false in case the verification fails.
func (t *T) CheckEqual(obtained interface{}, expected interface{},
                       issue ...interface{}) bool {
    summary := "CheckEqual(obtained, expected):"
    return t.internalCheckEqual(obtained, expected, true, summary, issue...)
}

// Verify if the first value is not equal to the second value.  In case
// they are equal, an error will be logged, the test will be marked as
// failed, and the test execution will continue.  The extra arguments are
// optional and, if provided, will be assembled together with fmt.Sprint()
// and printed next to the reported problem in case of errors. The returned
// value will be false in case the verification fails.
func (t *T) CheckNotEqual(obtained interface{}, expected interface{},
                          issue ...interface{}) bool {
    summary := "CheckNotEqual(obtained, unexpected):"
    return t.internalCheckEqual(obtained, expected, false, summary, issue...)
}

// Ensure that the first value is equal to the second value.  In case
// they're not equal, an error will be logged, the test will be marked as
// failed, and the test execution will stop.  The extra arguments are
// optional and, if provided, will be assembled together with fmt.Sprint()
// and printed next to the reported problem in case of errors.
func (t *T) AssertEqual(obtained interface{}, expected interface{},
                        issue ...interface{}) {
    summary := "AssertEqual(obtained, expected):"
    if !t.internalCheckEqual(obtained, expected, true, summary, issue...) {
        t.stopNow()
    }
}

// Ensure that the first value is not equal to the second value.  In case
// they are equal, an error will be logged, the test will be marked as
// failed, and the test execution will stop.  The extra arguments are
// optional and, if provided, will be assembled together with fmt.Sprint()
// and printed next to the reported problem in case of errors.
func (t *T) AssertNotEqual(obtained interface{}, expected interface{},
                           issue ...interface{}) {
    summary := "AssertNotEqual(obtained, unexpected):"
    if !t.internalCheckEqual(obtained, expected, false, summary, issue...) {
        t.stopNow()
    }
}


func (t *T) internalCheckEqual(a interface{}, b interface{}, equal bool,
                               summary string, issue ...interface{}) bool {
    typeA := reflect.Typeof(a)
    typeB := reflect.Typeof(b)
    if (typeA == typeB && checkEqual(a, b)) != equal {
        t.logCaller(2, summary)
        if equal {
            t.logValue("Obtained", a)
            t.logValue("Expected", b)
        } else {
            t.logValue("Both", a)
        }
        if len(issue) != 0 {
            t.logString(fmt.Sprint(issue...))
        }
        t.logNewLine()
        t.Fail()
        return false
    }
    return true
}

// This will use a fast path to check for equality of normal types,
// and then fallback to reflect.DeepEqual if things go wrong.
func checkEqual(a interface{}, b interface{}) (result bool) {
    defer func() {
        if recover() != nil {
            result = reflect.DeepEqual(a, b)
        }
    }()
    return (a == b)
}


// -----------------------------------------------------------------------
// String matching testing.

// Verify if the value provided matches with the given regular expression.
// The value must be either a string, or a value which provides the String()
// method. In case it doesn't match, an error will be logged, the test will
// be marked as failed, and the test execution will continue. The extra
// arguments are optional and, if provided, will be assembled together with
// fmt.Sprint() and printed next to the reported problem in case of errors.
func (t *T) CheckMatch(value interface{}, expression string,
                       issue ...interface{}) bool {
    summary := "CheckMatch(value, expression):"
    return t.internalCheckMatch(value, expression, true, summary, issue...)
}

// Ensure that the value provided matches with the given regular expression.
// The value must be either a string, or a value which provides the String()
// method. In case it doesn't match, an error will be logged, the test will
// be marked as failed, and the test execution will stop. The extra
// arguments are optional and, if provided, will be assembled together with
// fmt.Sprint() and printed next to the reported problem in case of errors.
func (t *T) AssertMatch(value interface{}, expression string,
                        issue ...interface{}) {
    summary := "AssertMatch(value, expression):"
    if !t.internalCheckMatch(value, expression, true, summary, issue...) {
        t.stopNow()
    }
}

func (t *T) internalCheckMatch(value interface{}, expression string,
                               equal bool, summary string,
                               issue ...interface{}) bool {
    valueStr, valueIsStr := value.(string)
    if !valueIsStr {
        if valueWithStr, valueHasStr := value.(hasString); valueHasStr {
            valueStr, valueIsStr = valueWithStr.String(), true
        }
    }
    var err os.Error
    var matches bool
    if valueIsStr {
        matches, err = regexp.MatchString("^" + expression + "$", valueStr)
    }
    if !matches || err != nil {
        t.logCaller(2, summary)
        var msg string
        if !matches {
            t.logValue("Value", value)
            msg = fmt.Sprintf("Expected to match expression: %#v", expression)
        } else {
            msg = fmt.Sprintf("Can't compile match expression: %#v", expression)
        }
        t.logString(msg)
        if len(issue) != 0 {
            t.logString(fmt.Sprint(issue...))
        }
        t.logNewLine()
        t.Fail()
        return false
    }
    return true
}
