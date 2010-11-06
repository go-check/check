package gocheck

import (
    "reflect"
    "runtime"
    "strings"
    "strconv"
    "regexp"
    "path"
    "sync"
    "rand"
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

type callKind int

const (
    succeededSt = iota
    failedSt
    skippedSt
    panickedSt
    fixturePanickedSt
    missedSt
)

type callStatus int

type call struct {
    method *reflect.FuncValue
    kind callKind
    status callStatus
    logv string
    done chan *call
    expectedFailure *string
    tempDir *tempDir
}

func newCall(method *reflect.FuncValue, kind callKind, tempDir *tempDir) *call {
    return &call{method:method, kind:kind, tempDir:tempDir,
                 done:make(chan *call, 1)}
}

func (c *call) stopNow() {
    runtime.Goexit()
}


// -----------------------------------------------------------------------
// XXX Where to put these?

// Type passed as an argument to fixture methods.
type F struct {
    *call
}


// -----------------------------------------------------------------------
// Handling of temporary files and directories.

type tempDir struct {
    sync.Mutex
    _path string
    _counter int
}

func (td *tempDir) newPath() string {
    td.Lock()
    defer td.Unlock()
    if td._path == "" {
        var err os.Error
        for i := 0; i != 100; i++ {
            path := fmt.Sprintf("%s/gocheck-%d", os.TempDir(), rand.Int())
            if err = os.Mkdir(path, 0700); err == nil {
                td._path = path
                break
            }
        }
        if td._path == "" {
            panic("Couldn't create temporary directory: " + err.String())
        }
    }
    result := path.Join(td._path, strconv.Itoa(td._counter))
    td._counter += 1
    return result
}

func (td *tempDir) removeAll() {
    td.Lock()
    defer td.Unlock()
    if td._path != "" {
        err := os.RemoveAll(td._path)
        if err != nil {
            println("WARNING: Error cleaning up temporaries: " + err.String())
        }
    }
}

func (c *call) MkDir() string {
    path := c.tempDir.newPath()
    if err := os.Mkdir(path, 0700); err != nil {
        panic(fmt.Sprintf("Couldn't create temporary directory %s: %s",
                          path, err.String()))
    }
    return path
}

// The methods below are not strictly necessary, but godoc doesn't yet
// understand the method above is public due to the embedded type.

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (t *T) MkDir() string {
    return t.call.MkDir()
}

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (f *F) MkDir() string {
    return f.call.MkDir()
}


// -----------------------------------------------------------------------
// Low-level logging functions.

func (c *call) log(args ...interface{}) {
    c.logv += fmt.Sprint(args...) + "\n"
}

func (c *call) logf(format string, args ...interface{}) {
    c.logv += fmt.Sprintf(format, args...) + "\n"
}

func (c *call) logNewLine() {
    c.logv += "\n"
}

type hasString interface {
    String() string
}

func (c *call) logValue(label string, value interface{}) {
    if label == "" {
        if _, ok := value.(hasString); ok {
            c.logf("... %#v (%v)", value, value)
        } else {
            c.logf("... %#v", value)
        }
    } else if (value == nil) {
        c.logf("... %s (nil): nil", label)
    } else {
        if _, ok := value.(hasString); ok {
            c.logf("... %s (%s): %#v (%v)",
                   label, reflect.Typeof(value).String(), value, value)
        } else {
            c.logf("... %s (%s): %#v",
                   label, reflect.Typeof(value).String(), value)
        }
    }
}

func (c *call) logString(issue string) {
    c.log("... ", issue)
}

func (c *call) logCaller(skip int, issue string) {
    // This is a bit heavier than it ought to be.
    skip += 1 // Our own frame.
    if pc, callerFile, callerLine, ok := runtime.Caller(skip); ok {
        var testFile string
        var testLine int
        testFunc := runtime.FuncForPC(c.method.Get())
        if runtime.FuncForPC(pc) != testFunc {
            for {
                skip += 1
                if pc, file, line, ok := runtime.Caller(skip); ok {
                    // Note that the test line may be different on
                    // distinct calls for the same test.  Showing
                    // the "internal" line is helpful when debugging.
                    if runtime.FuncForPC(pc) == testFunc {
                        testFile, testLine = file, line
                        break
                    }
                } else {
                    break
                }
            }
        }
        if testFile != "" && (testFile != callerFile ||
                              testLine != callerLine) {
            c.logf("%s:%d > %s:%d:\n... %s", nicePath(testFile), testLine,
                  nicePath(callerFile), callerLine, issue)
        } else {
            c.logf("%s:%d:\n... %s", nicePath(callerFile), callerLine, issue)
        }
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
        filename, line := function.FileLine(pc)
        return fmt.Sprintf("%s:%d", nicePath(filename), line)
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

type Result struct {
    Succeeded int
    Failed int
    Skipped int
    Panicked int
    FixturePanicked int
    Missed int // Not even tried to run, related to a panic in fixture.
    RunError os.Error // Houston, we've got a problem.
}

type resultTracker struct {
    writer io.Writer
    result Result
    _waiting int
    _missed int
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
                    handleExpectedFailure(c)
                    switch c.status {
                        case succeededSt:
                            if c.kind == testKd {
                                tracker.result.Succeeded++
                            }
                        case failedSt:
                            tracker.result.Failed++
                            tracker._reportProblem("FAIL", c)
                        case panickedSt:
                            if c.kind == fixtureKd {
                                tracker.result.FixturePanicked++
                            } else {
                                tracker.result.Panicked++
                            }
                            tracker._reportProblem("PANIC", c)
                        case fixturePanickedSt:
                            // That's a testKd call reporting that its fixture
                            // has panicked. The fixture call which caused the
                            // panic itself was tracked above. We'll report to
                            // aid debugging.
                            tracker._reportProblem("PANIC", c)
                            // And will track it as missed, since the panic
                            // was on the fixture, not on the test.
                            tracker.result.Missed++
                        case missedSt:
                            tracker.result.Missed++
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

func handleExpectedFailure(c *call) {
    if c.expectedFailure != nil {
        switch c.status {
            case failedSt:
                c.status = succeededSt
            case succeededSt:
                c.status = failedSt
                c.logString("Error: Test succeeded, but was expected to fail")
                c.logString("Reason: " + *c.expectedFailure)
        }
    }
}

func (tracker *resultTracker) _reportProblem(label string, c *call) {
    pc := c.method.Get()
    header := fmt.Sprintf(
        "\n-----------------------------------" +
        "-----------------------------------\n" +
        "%s: %s: %s\n\n",
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
    tempDir *tempDir
}

type RunConf struct {
    Output io.Writer
    Filter string
}

// Create a new suiteRunner able to run all methods in the given suite.
func newSuiteRunner(suite interface{}, runConf *RunConf) *suiteRunner {
    var output io.Writer
    var filter string

    output = os.Stdout

    if runConf != nil {
        if runConf.Output != nil {
            output = runConf.Output
        }
        if runConf.Filter != "" {
            filter = "^" + runConf.Filter + "$"
        }
    }

    suiteType := reflect.Typeof(suite)
    suiteNumMethods := suiteType.NumMethod()
    suiteValue := reflect.NewValue(suite)

    runner := suiteRunner{suite:suite,
                          tracker:newResultTracker(output)}
    runner.tests = make([]*reflect.FuncValue, suiteNumMethods)
    runner.tempDir = new(tempDir)
    testsLen := 0

    var filterRegexp *regexp.Regexp
    if filter != "" {
        if regexp, err := regexp.Compile(filter); err != nil {
            msg := "Bad filter expression: " + err.String()
            runner.tracker.result.RunError = os.NewError(msg)
            return &runner
        } else {
            filterRegexp = regexp
        }
    }

    // This map will be used to filter out duplicated methods.  This
    // looks like a bug in Go, described on issue 906:
    // http://code.google.com/p/go/issues/detail?id=906
    seen := make(map[uintptr]bool, suiteNumMethods)

    // XXX Shouldn't Name() work here? Why does it return an empty string?
    suiteName := suiteType.String()
    if index := strings.LastIndex(suiteName, "."); index != -1 {
        suiteName = suiteName[index+1:]
    }

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
                if isWantedTest(suiteName, method.Name, filterRegexp) {
                    runner.tests[testsLen] = funcValue
                    testsLen += 1
                }
        }
    }

    runner.tests = runner.tests[0:testsLen]
    return &runner
}

// Return true if the given suite name and method name should be
// considered as a test to be run.
func isWantedTest(suiteName, testName string, filterRegexp *regexp.Regexp) bool {
    if !strings.HasPrefix(testName, "Test") {
        return false
    } else if filterRegexp == nil {
        return true
    }
    return (filterRegexp.MatchString(testName) ||
            filterRegexp.MatchString(suiteName) ||
            filterRegexp.MatchString(suiteName + "." + testName))
}


// Run all methods in the given suite.
func (runner *suiteRunner) run() *Result {
    if runner.tracker.result.RunError == nil && len(runner.tests) > 0 {
        runner.tracker.start()
        if runner.checkFixtureArgs() {
            if runner.runFixture(runner.setUpSuite) {
                for i := 0; i != len(runner.tests); i++ {
                    c := runner.runTest(runner.tests[i])
                    if c.status == fixturePanickedSt {
                        runner.missTests(runner.tests[i+1:])
                        break
                    }
                }
            } else {
                runner.missTests(runner.tests)
            }
            runner.runFixture(runner.tearDownSuite)
        } else {
            runner.missTests(runner.tests)
        }
        runner.tracker.waitAndStop()
        runner.tempDir.removeAll()
    }
    return &runner.tracker.result
}


// Create a call object with the given suite method, and fork a
// goroutine with the provided dispatcher for running it.
func (runner *suiteRunner) forkCall(method *reflect.FuncValue, kind callKind,
                                    dispatcher func(c *call)) *call {
    c := newCall(method, kind, runner.tempDir)
    runner.tracker.waitForCall(c)
    go (func() {
        defer runner.callDone(c)
        dispatcher(c)
    })()
    return c
}

// Same as forkCall(), but wait for call to finish before returning.
func (runner *suiteRunner) runCall(method *reflect.FuncValue, kind callKind,
                                   dispatcher func(c *call)) *call {
    c := runner.forkCall(method, kind, dispatcher)
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
        c := runner.runCall(method, fixtureKd, func(c *call) {
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
    return runner.forkCall(method, testKd, func(c *call) {
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

// Helper to mark tests as missed.  A bit heavy for what it does, but it
// enables homogeneous handling of tracking, including nice verbose output.
func (runner *suiteRunner) missTests(methods []*reflect.FuncValue) {
    for _, method := range methods {
        runner.runCall(method, testKd, func(c *call) {
            c.status = missedSt
        })
    }
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
                runner.runCall(fv, fixtureKd, func(c *call) {
                    c.logArgPanic(fv, "*gocheck.F")
                    c.status = panickedSt
                })
            }
        }
    }
    return succeeded
}
