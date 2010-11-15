/*
Gocheck - A rich testing framework for Go

Copyright (c) 2010, Gustavo Niemeyer <gustavo@niemeyer.net>

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of the copyright holder nor the names of its
      contributors may be used to endorse or promote products derived from
      this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

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

type funcKind int

const (
    succeededSt = iota
    failedSt
    skippedSt
    panickedSt
    fixturePanickedSt
    missedSt
)

type funcStatus int

type C struct {
    method *reflect.FuncValue
    kind funcKind
    status funcStatus
    logv string
    done chan *C
    expectedFailure *string
    tempDir *tempDir
}

func newC(method *reflect.FuncValue, kind funcKind, tempDir *tempDir) *C {
    return &C{method:method, kind:kind, tempDir:tempDir, done:make(chan *C, 1)}
}

func (c *C) stopNow() {
    runtime.Goexit()
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

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (c *C) MkDir() string {
    path := c.tempDir.newPath()
    if err := os.Mkdir(path, 0700); err != nil {
        panic(fmt.Sprintf("Couldn't create temporary directory %s: %s",
                          path, err.String()))
    }
    return path
}


// -----------------------------------------------------------------------
// Low-level logging functions.

func (c *C) log(args ...interface{}) {
    c.logv += fmt.Sprint(args...) + "\n"
}

func (c *C) logf(format string, args ...interface{}) {
    c.logv += fmt.Sprintf(format, args...) + "\n"
}

func (c *C) logNewLine() {
    c.logv += "\n"
}

type hasString interface {
    String() string
}

func (c *C) logValue(label string, value interface{}) {
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

func (c *C) logString(issue string) {
    c.log("... ", issue)
}

func (c *C) logCaller(skip int, issue string) {
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

func (c *C) logPanic(skip int, value interface{}) {
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

func (c *C) logSoftPanic(issue string) {
    c.log("... Panic: ", issue)
}

func (c *C) logArgPanic(funcValue *reflect.FuncValue, expectedType string) {
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
    Missed int // Not even tried to run, related to a panic in the fixture.
    RunError os.Error // Houston, we've got a problem.
}

type resultTracker struct {
    writer io.Writer
    result Result
    _waiting int
    _missed int
    _waitChan chan *C
    _doneChan chan *C
    _stopChan chan bool
}

func newResultTracker(writer io.Writer) *resultTracker {
    return &resultTracker{writer: writer,
                          _waitChan: make(chan *C),     // Synchronous
                          _doneChan: make(chan *C, 32), // Asynchronous
                          _stopChan: make(chan bool)}      // Synchronous
}

func (tracker *resultTracker) start() {
    go tracker._loopRoutine()
}

func (tracker *resultTracker) waitAndStop() {
    <-tracker._stopChan
}

func (tracker *resultTracker) waitForCall(c *C) {
    tracker._waitChan <- c
}

func (tracker *resultTracker) callDone(c *C) {
    tracker._doneChan <- c
}

func (tracker *resultTracker) _loopRoutine() {
    for {
        var c *C
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

func handleExpectedFailure(c *C) {
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

func (tracker *resultTracker) _reportProblem(label string, c *C) {
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
            filter = runConf.Filter
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
func (runner *suiteRunner) forkCall(method *reflect.FuncValue, kind funcKind,
                                    dispatcher func(c *C)) *C {
    c := newC(method, kind, runner.tempDir)
    runner.tracker.waitForCall(c)
    go (func() {
        defer runner.callDone(c)
        dispatcher(c)
    })()
    return c
}

// Same as forkCall(), but wait for call to finish before returning.
func (runner *suiteRunner) runFunc(method *reflect.FuncValue, kind funcKind,
                                   dispatcher func(c *C)) *C {
    c := runner.forkCall(method, kind, dispatcher)
    <-c.done
    return c
}

// Handle a finished call.  If there were any panics, update the call status
// accordingly.  Then, mark the call as done and report to the tracker.
func (runner *suiteRunner) callDone(c *C) {
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
        c := runner.runFunc(method, fixtureKd, func(c *C) {
            c.method.Call([]reflect.Value{reflect.NewValue(c)})
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
func (runner *suiteRunner) forkTest(method *reflect.FuncValue) *C {
    return runner.forkCall(method, testKd, func(c *C) {
        defer runner.runFixtureWithPanic(runner.tearDownTest)
        runner.runFixtureWithPanic(runner.setUpTest)
        methodType := c.method.Type().(*reflect.FuncType)
        if methodType.In(1) == reflect.Typeof(c) && methodType.NumIn() == 2 {
            c.method.Call([]reflect.Value{reflect.NewValue(c)})
        } else {
            // Rather than a plain panic, provide a more helpful message when
            // the argument type is incorrect.
            c.status = panickedSt
            c.logArgPanic(c.method, "*gocheck.C")
        }
    })
}

// Same as forkTest(), but wait for the test to finish before returning.
func (runner *suiteRunner) runTest(method *reflect.FuncValue) *C {
    c := runner.forkTest(method)
    <-c.done
    return c
}

// Helper to mark tests as missed.  A bit heavy for what it does, but it
// enables homogeneous handling of tracking, including nice verbose output.
func (runner *suiteRunner) missTests(methods []*reflect.FuncValue) {
    for _, method := range methods {
        runner.runFunc(method, testKd, func(c *C) {
            c.status = missedSt
        })
    }
}

// Verify if the fixture arguments are *gocheck.C.  In case of errors,
// log the error as a panic in the fixture method call, and return false.
func (runner *suiteRunner) checkFixtureArgs() bool {
    succeeded := true
    argType := reflect.Typeof(&C{})
    for _, fv := range []*reflect.FuncValue{runner.setUpSuite,
                                            runner.tearDownSuite,
                                            runner.setUpTest,
                                            runner.tearDownTest} {
        if fv != nil {
            fvType := fv.Type().(*reflect.FuncType)
            if fvType.In(1) != argType || fvType.NumIn() != 2 {
                succeeded = false
                runner.runFunc(fv, fixtureKd, func(c *C) {
                    c.logArgPanic(fv, "*gocheck.C")
                    c.status = panickedSt
                })
            }
        }
    }
    return succeeded
}
