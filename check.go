// Package check is a rich testing extension for Go's testing package.
//
// For details about the project, see:
//
//     http://labix.org/gocheck
//
package check

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
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

type tempDir struct {
	sync.Mutex
	path    string
	counter int
}

func (td *tempDir) newPath() string {
	td.Lock()
	defer td.Unlock()
	if td.path == "" {
		path, err := ioutil.TempDir("", "check-")
		if err != nil {
			panic("Couldn't create temporary directory: " + err.Error())
		}
		td.path = path
	}
	result := filepath.Join(td.path, strconv.Itoa(td.counter))
	td.counter++
	return result
}

func (td *tempDir) removeAll() {
	td.Lock()
	defer td.Unlock()
	if td.path != "" {
		err := os.RemoveAll(td.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: Error cleaning up temporaries: "+err.Error())
		}
	}
}

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (c *C) MkDir() string {
	path := c.tempDir.newPath()
	if err := os.Mkdir(path, 0700); err != nil {
		panic(fmt.Sprintf("Couldn't create temporary directory %s: %s", path, err.Error()))
	}
	return path
}

// -----------------------------------------------------------------------
// Low-level logging functions.

func (c *C) logNewLine() {
	c.Helper()
	c.Log()
}


func hasStringOrError(x interface{}) (ok bool) {
	_, ok = x.(fmt.Stringer)
	if ok {
		return
	}
	_, ok = x.(error)
	return
}

func (c *C) logValue(label string, value interface{}) {
	c.Helper()
	if label == "" {
		if hasStringOrError(value) {
			c.Logf("... %#v (%q)", value, value)
		} else {
			c.Logf("... %#v", value)
		}
	} else if value == nil {
		c.Logf("... %s = nil", label)
	} else {
		if hasStringOrError(value) {
			fv := fmt.Sprintf("%#v", value)
			qv := fmt.Sprintf("%q", value)
			if fv != qv {
				c.Logf("... %s %s = %s (%s)", label, reflect.TypeOf(value), fv, qv)
				return
			}
		}
		if s, ok := value.(string); ok && isMultiLine(s) {
			c.Logf(`... %s %s = "" +`, label, reflect.TypeOf(value))
			c.logMultiLine(s)
		} else {
			c.Logf("... %s %s = %#v", label, reflect.TypeOf(value), value)
		}
	}
}

func formatMultiLine(s string, quote bool) []byte {
	b := make([]byte, 0, len(s)*2)
	i := 0
	n := len(s)
	for i < n {
		j := i + 1
		for j < n && s[j-1] != '\n' {
			j++
		}
		b = append(b, "...     "...)
		if quote {
			b = strconv.AppendQuote(b, s[i:j])
		} else {
			b = append(b, s[i:j]...)
			b = bytes.TrimSpace(b)
		}
		if quote && j < n {
			b = append(b, " +"...)
		}
		b = append(b, '\n')
		i = j
	}
	return b
}

func (c *C) logMultiLine(s string) {
	c.Helper()
	c.Log(formatMultiLine(s, true))
}

func isMultiLine(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '\n' {
			return true
		}
	}
	return false
}

func (c *C) logString(issue string) {
	c.Helper()
	c.Log("... ", issue)
}

func (c *C) logCaller(skip int) {
	c.Helper()
	// This is a bit heavier than it ought to be.
	skip++ // Our own frame.
	pc, callerFile, callerLine, ok := runtime.Caller(skip)
	if !ok {
		return
	}
	var testFile string
	var testLine int
	testFunc := runtime.FuncForPC(c.method.PC())
	if runtime.FuncForPC(pc) != testFunc {
		for {
			skip++
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
	if testFile != "" && (testFile != callerFile || testLine != callerLine) {
		c.logCode(testFile, testLine)
	}
	c.logCode(callerFile, callerLine)
}

func (c *C) logCode(path string, line int) {
	c.Helper()
	c.Logf("%s:%d:", nicePath(path), line)
	code, err := printLine(path, line)
	if code == "" {
		code = "..." // XXX Open the file and take the raw line.
		if err != nil {
			code += err.Error()
		}
	}
	c.Log(indent(code, "    "))
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

type RunConf struct {
	KeepWorkDir bool
}

// Create a new suiteRunner able to run all methods in the given suite.
func newSuiteRunner(suite interface{}, runConf *RunConf) *suiteRunner {
	var conf RunConf
	if runConf != nil {
		conf = *runConf
	}

	suiteType := reflect.TypeOf(suite)
	suiteNumMethods := suiteType.NumMethod()
	suiteValue := reflect.ValueOf(suite)

	runner := &suiteRunner{
		suite:   suite,
		tempDir: &tempDir{},
		keepDir: conf.KeepWorkDir,
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
		if runner.keepDir {
			t.Log("WORK =", runner.tempDir.path)
		} else {
			runner.tempDir.removeAll()
		}
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
