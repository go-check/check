package gocheck


import (
    "fmt"
    "reflect"
)


type T struct {
    *call
}


// -----------------------------------------------------------------------
// Basic succeeding/failing logic.

func (t *T) Failed() bool {
    return t.call.status == failedSt
}

func (t *T) Fail() {
    t.call.status = failedSt
}

func (t *T) FailNow() {
    t.Fail()
    t.stopNow()
}

func (t *T) Succeed() {
    t.call.status = succeededSt
}

func (t *T) SucceedNow() {
    t.Succeed()
    t.stopNow()
}

func (t *T) ExpectFailure(reason string) {
    t.expectedFailure = &reason
}


// -----------------------------------------------------------------------
// Basic logging.

func (t *T) GetLog() string {
    return t.call.logv
}

func (t *T) Log(args ...interface{}) {
    t.log(args...)
}

func (t *T) Logf(format string, args ...interface{}) {
    t.logf(format, args...)
}

func (t *T) Error(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    t.Fail()
}

func (t *T) Errorf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprintf("Error: " + format, args...))
    t.Fail()
}

func(t *T) Fatal(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args...)))
    t.FailNow()
}

func(t *T) Fatalf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprintf(format, args...)))
    t.FailNow()
}


// -----------------------------------------------------------------------
// Testing helper functions.

func (t *T) CheckEqual(obtained interface{}, expected interface{},
                       issue ...interface{}) bool {
    return t.internalCheckEqual(obtained, expected, true,
                                "CheckEqual(A, B): A != B", issue...)
}

func (t *T) CheckNotEqual(obtained interface{}, expected interface{},
                          issue ...interface{}) bool {
    return t.internalCheckEqual(obtained, expected, false,
                                "CheckNotEqual(A, B): A == B", issue...)
}

func (t *T) AssertEqual(obtained interface{}, expected interface{},
                        issue ...interface{}) {
    if !t.internalCheckEqual(obtained, expected, true,
                             "AssertEqual(A, B): A != B", issue...) {
        t.stopNow()
    }
}

func (t *T) AssertNotEqual(obtained interface{}, expected interface{},
                           issue ...interface{}) {
    if !t.internalCheckEqual(obtained, expected, false,
                             "AssertNotEqual(A, B): A == B", issue...) {
        t.stopNow()
    }
}


func (t *T) internalCheckEqual(a interface{}, b interface{}, equal bool,
                               summary string, issue ...interface{}) bool {
    if checkEqual(a, b) != equal {
        t.logCaller(2, summary)
        t.logValue("A:", a)
        t.logValue("B:", b)
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
