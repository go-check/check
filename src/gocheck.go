package gocheck

import (
    "reflect"
    "runtime"
    "strings"
    "path"
    "fmt"
    "io"
    "os"
)


// -----------------------------------------------------------------------
// Internal type which deals with suite method calling.

const (
    fixtureKd = iota
    testKd
    benchmarkKd
)

const (
    succeededSt = iota
    failedSt
    skippedSt
    panickedSt
    fixturePanickedSt
)

type call struct {
    method *reflect.FuncValue
    kind int
    status int
    logv string
    done chan *call
}

func newCall(method *reflect.FuncValue) *call {
    return &call{method:method, done:make(chan *call, 1)}
}

func (c *call) stopNow() {
    runtime.Goexit()
}


// -----------------------------------------------------------------------
// XXX Where to put these?

type F struct {
    *call
}


// -----------------------------------------------------------------------
// Low-level logging functions.

func (c *call) log(args ...interface{}) {
    c.logv += fmt.Sprint(args) + "\n"
}

func (c *call) logf(format string, args ...interface{}) {
    c.logv += fmt.Sprintf(format, args) + "\n"
}

func (c *call) logNewLine() {
    c.logv += "\n"
}

func (c *call) logValue(label string, value interface{}) {
    if label == "" {
        c.logf("... %#v", value)
    } else {
        c.logf("... %s %#v", label, value)
    }
}

func (c *call) logString(issue string) {
    c.log("... ", issue)
}

func (c *call) logCaller(skip int, issue string) {
    if _, callerFile, callerLine, ok := runtime.Caller(skip+1); ok {
        c.logf("%s:%d:\n... %s", nicePath(callerFile), callerLine, issue)
    }
}

func (c *call) logPanic(skip int, value interface{}) {
    skip += 1 // Our own frame.
    initialSkip := skip
    for {
        if pc, file, line, ok := runtime.Caller(skip); ok {
            if skip == initialSkip {
                c.logf("... Panic: %s (PC=0x%X)\n", value, pc)
            }
            name := niceFuncName(pc)
            if name == "reflect.FuncValue.Call" ||
               name == "gocheck.forkTest" {
                break
            }
            c.logf("%s:%d\n  in %s", nicePath(file), line, name)
        } else {
            break
        }
        skip += 1
    }
}

func (c *call) logSoftPanic(issue string) {
    c.log("... Panic: ", issue)
}

func (c *call) logArgPanic(funcValue *reflect.FuncValue, expectedType string) {
    c.logf("... Panic: %s argument should be %s",
           niceFuncName(funcValue.Get()), expectedType)
}


// -----------------------------------------------------------------------
// Some simple formatting helpers.

var initWD, initWDErr = os.Getwd()

func nicePath(path string) string {
    if initWDErr == nil {
        if strings.HasPrefix(path, initWD+"/") {
            return path[len(initWD)+1:]
        }
    }
    return path
}

func niceFuncPath(pc uintptr) string {
    function := runtime.FuncForPC(pc)
    if function != nil {
        filename, _ := function.FileLine(pc)
        return nicePath(filename)
    }
    return "<unknown path>"
}

func niceFuncName(pc uintptr) string {
    function := runtime.FuncForPC(pc)
    if function != nil {
        name := path.Base(function.Name())
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
    return "<unknown function>"
}


// -----------------------------------------------------------------------
// Result tracker to aggregate call results.

type resultTracker struct {
    writer io.Writer
    _waiting int
    _waitChan chan *call
    _doneChan chan *call
    _stopChan chan bool
}

func newResultTracker(writer io.Writer) *resultTracker {
    return &resultTracker{writer: writer,
                          _waitChan: make(chan *call),     // Synchronous
                          _doneChan: make(chan *call, 32), // Asynchronous
                          _stopChan: make(chan bool)}      // Synchronous
}

func (tracker *resultTracker) start() {
    go tracker._loopRoutine()
}

func (tracker *resultTracker) waitAndStop() {
    <-tracker._stopChan
}

func (tracker *resultTracker) waitForCall(c *call) {
    tracker._waitChan <- c
}

func (tracker *resultTracker) callDone(c *call) {
    tracker._doneChan <- c
}

func (tracker *resultTracker) _loopRoutine() {
    for {
        var c *call
        if tracker._waiting > 0 {
            // Calls still running. Can't stop.
            select {
                case c = <-tracker._waitChan:
                    tracker._waiting += 1
                case c = <-tracker._doneChan:
                    tracker._waiting -= 1
                    switch c.status {
                        case failedSt:
                            tracker._reportProblem("FAIL", c)
                        case panickedSt:
                            tracker._reportProblem("PANIC", c)
                        case fixturePanickedSt:
                            tracker._reportProblem("PANIC", c)
                    }
            }
        } else {
            // No calls.  Can stop, but no done calls here.
            select {
                case tracker._stopChan <- true:
                    return
                case c = <-tracker._waitChan:
                    tracker._waiting += 1
                case c = <-tracker._doneChan:
                    panic("Tracker got an unexpected done call.")
            }
        }
    }
}

func (tracker *resultTracker) _reportProblem(label string, c *call) {
    // XXX How to get the first line where a function was defined?
    pc := c.method.Get()
    header := fmt.Sprintf(
        "\n-----------------------------------" +
        "-----------------------------------\n" +
        "%s: %s:%s\n\n",
        label, niceFuncPath(pc), niceFuncName(pc))
    io.WriteString(tracker.writer, header)
    io.WriteString(tracker.writer, c.logv)
}



// -----------------------------------------------------------------------
// The underlying suite runner.

type suiteRunner struct {
    suite interface{}
    setUpSuite, tearDownSuite *reflect.FuncValue
    setUpTest, tearDownTest *reflect.FuncValue
    tests []*reflect.FuncValue
    tracker *resultTracker
}

// Create a new suiteRunner able to run all methods in the given suite.
func newSuiteRunner(suite interface{}, writer io.Writer) *suiteRunner {
    suiteType := reflect.Typeof(suite)
    suiteNumMethods := suiteType.NumMethod()
    suiteValue := reflect.NewValue(suite)

    runner := suiteRunner{suite:suite, tracker:newResultTracker(writer)}
    runner.tests = make([]*reflect.FuncValue, suiteNumMethods)
    testsLen := 0

    // This map will be used to filter out duplicated methods.  This
    // looks like a bug in Go, described on issue 906:
    // http://code.google.com/p/go/issues/detail?id=906
    seen := make(map[uintptr]bool, suiteNumMethods)

    for i := 0; i != suiteNumMethods; i++ {
        funcValue := suiteValue.Method(i)
        funcPC := funcValue.Get()
        if _, found := seen[funcPC]; found {
            continue
        }
        seen[funcPC] = true
        method := suiteType.Method(i)
        switch method.Name {
            case "SetUpSuite":
                runner.setUpSuite = funcValue
            case "TearDownSuite":
                runner.tearDownSuite = funcValue
            case "SetUpTest":
                runner.setUpTest = funcValue
            case "TearDownTest":
                runner.tearDownTest = funcValue
            default:
                if strings.HasPrefix(method.Name, "Test") {
                    runner.tests[testsLen] = funcValue
                    testsLen += 1
                }
        }
    }

    runner.tests = runner.tests[0:testsLen]
    return &runner
}

// Run all methods in the given suite.
func (runner *suiteRunner) run() {
    runner.tracker.start()
    if runner.checkFixtureArgs() {
        if runner.runFixture(runner.setUpSuite) {
            for i := 0; i != len(runner.tests); i++ {
                c := runner.runTest(runner.tests[i])
                if c.status == fixturePanickedSt {
                    // XXX Should count the tests not run as skipped.
                    break
                }
            }
        }
        runner.runFixture(runner.tearDownSuite)
    } else {
        // XXX Should mark tests as skipped here.
    }
    runner.tracker.waitAndStop()
}


// Create a call object with the given suite method, and fork a
// goroutine with the provided dispatcher for running it.
func (runner *suiteRunner) forkCall(method *reflect.FuncValue,
                                    dispatcher func(c *call)) *call {
    c := newCall(method)
    runner.tracker.waitForCall(c)
    go (func() {
        defer runner.callDone(c)
        dispatcher(c)
    })()
    return c
}

// Same as forkCall(), but wait for call to finish before returning.
func (runner *suiteRunner) runCall(method *reflect.FuncValue,
                                   dispatcher func(c *call)) *call {
    c := runner.forkCall(method, dispatcher)
    <-c.done
    return c
}

// Handle a finished call.  If there were any panics, update the call status
// accordingly.  Then, mark the call as done and report to the tracker.
func (runner *suiteRunner) callDone(c *call) {
    value := recover()
    if value != nil {
        switch v := value.(type) {
            case *fixturePanic:
                c.logSoftPanic("Fixture has panicked (see related PANIC)")
                c.status = fixturePanickedSt
            default:
                c.logPanic(1, value)
                c.status = panickedSt
        }
    }
    runner.tracker.callDone(c)
    c.done <- c
}

// Runs a fixture call synchronously.  The fixture will still be run in a
// goroutine like all suite methods, but this method will not return
// while the fixture goroutine is not done, because the fixture must be
// run in a desired order.
func (runner *suiteRunner) runFixture(method *reflect.FuncValue) bool {
    if method != nil {
        c := runner.runCall(method, func(c *call) {
            c.method.Call([]reflect.Value{reflect.NewValue(&F{c})})
        })
        return (c.status == succeededSt)
    }
    return true
}

// Run the fixture method with runFixture(), but panic with a fixturePanic{}
// in case the fixture method panics.  This makes it easier to track the
// fixture panic together with other call panics within forkTest().
func (runner *suiteRunner) runFixtureWithPanic(method *reflect.FuncValue) {
    if !runner.runFixture(method) {
        panic(&fixturePanic{method})
    }
}

type fixturePanic struct {
    method *reflect.FuncValue
}

// Run the suite test method, together with the test-specific fixture,
// asynchronously.
func (runner *suiteRunner) forkTest(method *reflect.FuncValue) *call {
    return runner.forkCall(method, func(c *call) {
        defer runner.runFixtureWithPanic(runner.tearDownTest)
        runner.runFixtureWithPanic(runner.setUpTest)
        t := &T{c}
        tt := t.method.Type().(*reflect.FuncType)
        if tt.In(1) == reflect.Typeof(t) && tt.NumIn() == 2 {
            t.method.Call([]reflect.Value{reflect.NewValue(t)})
        } else {
            // Rather than a plain panic, provide a more helpful message when
            // the argument type is incorrect.
            t.status = panickedSt
            t.logArgPanic(t.method, "*gocheck.T")
        }
    })
}

// Same as forkTest(), but wait for the test to finish before returning.
func (runner *suiteRunner) runTest(method *reflect.FuncValue) *call {
    c := runner.forkTest(method)
    <-c.done
    return c
}

// Verify if the fixture arguments are *gocheck.F.  In case of errors,
// log the error as a panic in the fixture method call, and return false.
func (runner *suiteRunner) checkFixtureArgs() bool {
    succeeded := true
    argType := reflect.Typeof(&F{})
    for _, fv := range []*reflect.FuncValue{runner.setUpSuite,
                                            runner.tearDownSuite,
                                            runner.setUpTest,
                                            runner.tearDownTest} {
        if fv != nil {
            fvt := fv.Type().(*reflect.FuncType)
            if fvt.In(1) != argType || fvt.NumIn() != 2 {
                succeeded = false
                runner.runCall(fv, func(c *call) {
                    c.logArgPanic(fv, "*gocheck.F")
                    c.status = panickedSt
                })
            }
        }
    }
    return succeeded
}
