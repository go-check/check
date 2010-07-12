package gocheck


import (
    "fmt"
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


// -----------------------------------------------------------------------
// Basic logging.

func (t *T) GetLog() string {
    return t.call.logv
}

func (t *T) Log(args ...interface{}) {
    t.log(args)
}

func (t *T) Logf(format string, args ...interface{}) {
    t.logf(format, args)
}

func (t *T) Error(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args)))
    t.Fail()
}

func (t *T) Errorf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprintf("Error: " + format, args))
    t.Fail()
}

func(t *T) Fatal(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args)))
    t.FailNow()
}

func(t *T) Fatalf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprintf(format, args)))
    t.FailNow()
}


// -----------------------------------------------------------------------
// Testing helper functions.

func (t *T) CheckEqual(expected interface{}, obtained interface{},
                       issue ...interface{}) bool {
    return t.internalCheckEqual(expected, obtained, true,
                                "CheckEqual(A, B): A != B", issue)
}

func (t *T) CheckNotEqual(expected interface{}, obtained interface{},
                          issue ...interface{}) bool {
    return t.internalCheckEqual(expected, obtained, false,
                                "CheckNotEqual(A, B): A == B", issue)
}

func (t *T) AssertEqual(expected interface{}, obtained interface{},
                        issue ...interface{}) {
    if !t.internalCheckEqual(expected, obtained, true,
                             "AssertEqual(A, B): A != B", issue) {
        t.stopNow()
    }
}

func (t *T) AssertNotEqual(expected interface{}, obtained interface{},
                           issue ...interface{}) {
    if !t.internalCheckEqual(expected, obtained, false,
                             "AssertNotEqual(A, B): A == B", issue) {
        t.stopNow()
    }
}


func (t *T) internalCheckEqual(a interface{}, b interface{}, equal bool,
                               summary string, issue ...interface{}) bool {
    if (a == b) != equal {
        t.logCaller(2, summary)
        t.logValue("A:", a)
        t.logValue("B:", b)
        if len(issue) != 0 {
            t.logString(fmt.Sprint(issue))
        }
        t.logNewLine()
        t.Fail()
        return false
    }
    return true
}
