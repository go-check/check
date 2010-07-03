package gocheck

import (
    "reflect"
    "runtime"
    "testing"
    "strings"
    "path"
    "fmt"
    "io"
    "os"
)

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

func (t *T) GetLog() string {
    return t.log
}

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

func(t *T) Fatal(args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprint(args)))
    t.FailNow()
}

func(t *T) Fatalf(format string, args ...interface{}) {
    t.logCaller(1, fmt.Sprint("Error: ", fmt.Sprintf(format, args)))
    t.FailNow()
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
    if _, callerFile, callerLine, ok := runtime.Caller(skip+1); ok {
        t.Logf("%s:%d:\n... %s", nicePath(callerFile), callerLine, issue)
    }
}

func (t *T) logPanic(skip int, value interface{}) {
    skip += 1 // Our own frame.
    initialSkip := skip
    for {
        if pc, file, line, ok := runtime.Caller(skip); ok {
            if skip == initialSkip {
                t.Logf("... Panic: %s (PC=0x%X)\n", value, pc)
            }
            var name string
            if f := runtime.FuncForPC(pc); f != nil {
                name = niceFuncName(f.Name())
                if name == "reflect.FuncValue.Call" {
                    break
                }
            } else {
                name = "<unknown function>"
            }
            t.Logf("%s:%d\n  in %s", nicePath(file), line, name)
        } else {
            break
        }
        skip += 1
    }
}

var initWD, initWDErr = os.Getwd()

func nicePath(path string) string {
    if initWDErr == nil {
        if strings.HasPrefix(path, initWD+"/") {
            return path[len(initWD)+1:]
        }
    }
    return path
}

func niceFuncName(name string) string {
    name = path.Base(name)
    if strings.HasPrefix(name, "_xtest_.*") {
        name = name[9:]
    }
    if i := strings.LastIndex(name, ".*"); i != -1 {
        name = name[0:i] + "." + name[i+2:]
    }
    if i := strings.LastIndex(name, "·"); i != -1 {
        name = name[0:i] + "." + name[i+2:]
    }
    return name
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


// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface{}

func Suite(suite interface{}) interface{} {
    return Suites(suite)[0]
}

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
// Test running logic.

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

func TestingT(testingT *testing.T) {
    RunAll()
}

func RunAll() {
    for _, suite := range allSuites {
        Run(suite)
    }
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
    value := recover()
    if value != nil {
        t.logPanic(1, value)
        t.Fail()
    }
    t.exit <- t
}

func writeFailure(t *T, writer io.Writer) {
    testLocation := ""
    testFuncName := ""
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
        testFuncName = niceFuncName(testFunc.Name())
    } else {
        testFuncName = t.method.Name
    }
    header := fmt.Sprintf(
        "\n-----------------------------------" +
        "-----------------------------------\n" +
        "FAIL: %s%s\n\n", testLocation, testFuncName)
    io.WriteString(writer, header)
    io.WriteString(writer, t.log)
}
