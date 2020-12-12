// Package check is a rich testing extension for Go's testing package.
//
// For details about the project, see:
//
//     http://labix.org/gocheck
//
package check

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

// -----------------------------------------------------------------------
// Internal type which deals with suite method calling.

const (
	fixtureKd = iota
	testKd
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

type funcStatus uint32

// A method value can't reach its own Method structure.
type methodType struct {
	reflect.Value
	Info reflect.Method
}

func newMethod(receiver reflect.Value, i int) *methodType {
	return &methodType{receiver.Method(i), receiver.Type().Method(i)}
}

func (method *methodType) PC() uintptr {
	return method.Info.Func.Pointer()
}

func (method *methodType) suiteName() string {
	t := method.Info.Type.In(0)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func (method *methodType) String() string {
	return method.suiteName() + "." + method.Info.Name
}

func (m *methodType) Call(c *C) {
	if m == nil {
		return
	}

	c.method = m
	defer func() {
		c.method = nil
	}()

	c.Helper()
	switch f := m.Interface().(type) {
	case func(*C):
		f(c)
	case func(*testing.T):
		f(c.T)
	default:
		c.Fatalf("bad signature for method %s: %T", m.Info.Name, m.Interface())
	}
}

func (method *methodType) matches(re *regexp.Regexp) bool {
	return (re.MatchString(method.Info.Name) ||
		re.MatchString(method.suiteName()) ||
		re.MatchString(method.String()))
}

type C struct {
	*testing.T

	method    *methodType
	reason    string
	mustFail  bool
	tempDir   *tempDir
	startTime time.Time
}

// -----------------------------------------------------------------------
// Handling of temporary files and directories.

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (c *C) MkDir() string {
	c.Helper()
	return c.TempDir()
}

// -----------------------------------------------------------------------
// Some simple formatting helpers.

var initWD, initWDErr = os.Getwd()

func init() {
	if initWDErr == nil {
		initWD = strings.Replace(initWD, "\\", "/", -1) + "/"
	}
}

func nicePath(path string) string {
	if initWDErr == nil {
		if strings.HasPrefix(path, initWD) {
			return path[len(initWD):]
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
		if i := strings.Index(name, "."); i > 0 {
			name = name[i+1:]
		}
		if strings.HasPrefix(name, "(*") {
			if i := strings.Index(name, ")"); i > 0 {
				name = name[2:i] + name[i+1:]
			}
		}
		if i := strings.LastIndex(name, ".*"); i != -1 {
			name = name[:i] + "." + name[i+2:]
		}
		if i := strings.LastIndex(name, "·"); i != -1 {
			name = name[:i] + "." + name[i+2:]
		}
		return name
	}
	return "<unknown function>"
}

func suiteName(suite interface{}) string {
	suiteType := reflect.TypeOf(suite)
	if suiteType.Kind() == reflect.Ptr {
		return suiteType.Elem().Name()
	}
	return suiteType.Name()
}

// -----------------------------------------------------------------------
// The underlying suite runner.

type suiteRunner struct {
	suite                     interface{}
	setUpSuite, tearDownSuite *methodType
	setUpTest, tearDownTest   *methodType
	tests                     []*methodType
	tempDir                   *tempDir
	keepDir                   bool
}

// Create a new suiteRunner able to run all methods in the given suite.
func newSuiteRunner(suite interface{}) *suiteRunner {
	suiteType := reflect.TypeOf(suite)
	suiteNumMethods := suiteType.NumMethod()
	suiteValue := reflect.ValueOf(suite)

	runner := &suiteRunner{
		suite:   suite,
		tempDir: &tempDir{},
		tests:   make([]*methodType, 0, suiteNumMethods),
	}

	for i := 0; i != suiteNumMethods; i++ {
		method := newMethod(suiteValue, i)
		switch method.Info.Name {
		case "SetUpSuite":
			runner.setUpSuite = method
		case "TearDownSuite":
			runner.tearDownSuite = method
		case "SetUpTest":
			runner.setUpTest = method
		case "TearDownTest":
			runner.tearDownTest = method
		default:
			if !strings.HasPrefix(method.Info.Name, "Test") {
				continue
			}
			runner.tests = append(runner.tests, method)
		}
	}
	return runner
}

// Run all methods in the given suite.
func (runner *suiteRunner) run(t *testing.T) {
	t.Cleanup(func() {
		runner.tempDir.removeAll()
	})

	c := C{T: t, startTime: time.Now()}
	t.Cleanup(func() { runner.tearDownSuite.Call(&c) })
	runner.setUpSuite.Call(&c)

	for _, test := range runner.tests {
		t.Run(test.Info.Name, func(t *testing.T) {
			runner.runTest(t, test)
		})
	}
}

// Same as forkTest(), but wait for the test to finish before returning.
func (runner *suiteRunner) runTest(t *testing.T, method *methodType) {
	c := C{T: t, startTime: time.Now()}

	t.Cleanup(func() { runner.tearDownTest.Call(&c) })
	runner.setUpTest.Call(&c)
	method.Call(&c)
}
