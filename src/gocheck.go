package gocheck

import (
    "reflect"
    "runtime"
    "testing"
    "strings"
    "fmt"
    "io"
    "os"
    //"flag"
)


//var blah *string = flag.String("blah", "Blah", "help for blah")

type T struct {
    suite interface{}
    method reflect.Method
    exit chan *T
    failed bool
    log string
}

type Result struct {
    failures int
}


var testHook func(*T)

func SetTestHook(newTestHook func(*T)) {
    testHook = newTestHook
}


func (t *T) GetLog() string {
    return t.log
}

// -----------------------------------------------------------------------
// Basic succeeding/failing logic.

func (t *T) Failed() bool {
    return t.failed
}

func (t *T) Fail() {
    t.failed = true
}

func (t *T) FailNow() {
    t.Fail()
    t.stopNow()
}

func (t *T) Succeed() {
    t.failed = false
}

func (t *T) SucceedNow() {
    t.Succeed()
    t.stopNow()
}

// This doesn't do much at the moment, but all stopping should go through
// here, so it's a useful control point.
func (t *T) stopNow() {
    runtime.Goexit()
}


// -----------------------------------------------------------------------
// Basic logging.

func (t *T) Log(args ...interface{}) {
    log := fmt.Sprint(args) + "\n"
    t.log += log
}

func (t *T) Logf(format string, args ...interface{}) {
    log := fmt.Sprintf(format, args) + "\n"
    t.log += log
}

func (t *T) Error(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args)))
    t.Fail()
}

func (t *T) Errorf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprintf("Error: " + format, args))
    t.Fail()
}


// -----------------------------------------------------------------------
// Internal logging helpers.

func (t *T) logNewLine() {
    t.log += "\n"
}

func (t *T) logValue(label string, value interface{}) {
    if label == "" {
        t.Logf("... %#v", value)
    } else {
        t.Logf("... %s %#v", label, value)
    }
}

func (t *T) logString(issue string) {
    t.Log("... ", issue)
}

func (t *T) logCaller(skip int, issue string) {
    t.Logf("... %s%s", t.formatCaller(skip+1), issue)
}

func (t *T) formatCaller(skip int) string {
    // That's not very well tested.  How to simulate a situation where
    // we can't get a caller or a function out of a PC?
    if _, callerFile, callerLine, ok := runtime.Caller(skip+1); ok {
        testPC := t.method.Func.Get()
        testFunc := runtime.FuncForPC(testPC)
        if testFunc == nil {
            return fmt.Sprintf("%d:", callerLine)
        } else {
            testFile, _ := testFunc.FileLine(testFunc.Entry())
            if testFile != callerFile {
                if wd, err := os.Getwd(); err == nil {
                    if strings.HasPrefix(callerFile, wd) {
                        callerFile = callerFile[len(wd)+1:]
                    }
                }
                return fmt.Sprintf("%s:%d:", callerFile, callerLine)
            } else {
                return fmt.Sprintf("%d:", callerLine)
            }
        }
    }
    return ""
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

func (t *T) internalCheckEqual(a interface{}, b interface{}, equal bool,
                               summary string, issue ...interface{}) bool {
    if (a == b) != equal {
        t.logNewLine()
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




type hasSetUpSuite interface {
    SetUpSuite()
}

type hasTearDownSuite interface {
    TearDownSuite()
}

type hasSetUpTest interface {
    SetUpTest()
}

type hasTearDownTest interface {
    TearDownTest()
}


//func (t *T) Run(suite interface{}) {
//    // No difference right now.  In the future, it'll do result aggregation.
//    Run(suite)
//}

func RunTestingT(suite interface{}, testingT *testing.T) {
    Run(suite)
}

func Run(suite interface{}) {
    RunWithWriter(suite, os.Stdout)
}

func RunWithWriter(suite interface{}, writer io.Writer) {
    suiteType := reflect.Typeof(suite)
    suiteNumMethods := suiteType.NumMethod()

    tests := make([]reflect.Method, suiteNumMethods)
    testCount := 0

    for i := 0; i != suiteNumMethods; i++ {
        method := suiteType.Method(i)
        if strings.HasPrefix(method.Name, "Test") {
            tests[testCount] = method
            testCount += 1
        }
    }


    if s, ok := suite.(hasSetUpSuite); ok {
        s.SetUpSuite()
    }

    for i := 0; i != testCount; i++ {
        t := T{suite:suite, method:tests[i], exit:make(chan *T)}
        go runTest(&t)
        <-t.exit

        if t.failed {
            writeFailure(&t, writer)
        }
    }

    if s, ok := suite.(hasTearDownSuite); ok {
        s.TearDownSuite()
    }
}

func runTest(t *T) {
    // TODO Check Func prototype before calling it.
    if s, ok := t.suite.(hasSetUpTest); ok {
        s.SetUpTest()
    }
    if s, ok := t.suite.(hasTearDownTest); ok {
        defer s.TearDownTest()
    }
    // XXX This is out of order.  We'll move this up as we add more tests.
    defer handleTestExit(t)
    t.method.Func.Call([]reflect.Value{reflect.NewValue(t.suite),
                                       reflect.NewValue(t)})
}

func handleTestExit(t *T) {
    // Do nothing with panics for now.
    recover()
    t.exit <- t
}

func writeFailure(t *T, writer io.Writer) {
    testLocation := ""
    testPC := t.method.Func.Get()
    testFunc := runtime.FuncForPC(testPC)
    if testFunc != nil {
        testFile, _ := testFunc.FileLine(testPC)
        if wd, err := os.Getwd(); err == nil {
            if strings.HasPrefix(testFile, wd) {
                testFile = testFile[len(wd)+1:]
            }
        }
        // XXX How to get the first line where a function was defined?
        //testLocation = fmt.Sprintf("%s:%d:", testFile, testLine)
        testLocation = fmt.Sprintf("%s:", testFile)
    }
    header := fmt.Sprintf(
        "-----------------------------------" +
        "-----------------------------------\n" +
        "FAIL: %s%s\n", testLocation, t.method.Name)
    io.WriteString(writer, header)
    io.WriteString(writer, t.log)
}
