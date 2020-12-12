package check

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"strconv"
	"testing"
	"time"
)

// TestName returns the current test name in the form "SuiteName.TestName"
func (c *C) TestName() string {
	return c.Name()
}

// -----------------------------------------------------------------------
// Basic logging.

// Output enables *C to be used as a logger in functions that require only
// the minimum interface of *log.Logger.
func (c *C) Output(calldepth int, s string) error {
	d := time.Now().Sub(c.startTime)
	msec := d / time.Millisecond
	sec := d / time.Second
	min := d / time.Minute

	c.Logf("[LOG] %d:%02d.%03d %s", min, sec%60, msec%1000, s)
	return nil
}

// -----------------------------------------------------------------------
// Generic checks and assertions based on checkers.

// Check verifies if the first value matches the expected value according
// to the provided checker. If they do not match, an error is logged, the
// test is marked as failed, and the test execution continues.
//
// Some checkers may not need the expected argument (e.g. IsNil).
//
// If the last value in args implements CommentInterface, it is used to log
// additional information instead of being passed to the checker (see Commentf
// for an example).
func (c *C) Check(obtained interface{}, checker Checker, args ...interface{}) bool {
	c.Helper()
	return internalCheck(c, "Check", obtained, checker, args...)
}

// Assert ensures that the first value matches the expected value according
// to the provided checker. If they do not match, an error is logged, the
// test is marked as failed, and the test execution stops.
//
// Some checkers may not need the expected argument (e.g. IsNil).
//
// If the last value in args implements CommentInterface, it is used to log
// additional information instead of being passed to the checker (see Commentf
// for an example).
func (c *C) Assert(obtained interface{}, checker Checker, args ...interface{}) {
	c.Helper()
	if !internalCheck(c, "Assert", obtained, checker, args...) {
		c.FailNow()
	}
}

func Check(c testing.TB, obtained interface{}, checker Checker, args ...interface{}) bool {
	c.Helper()
	return internalCheck(c, "Check", obtained, checker, args...)
}

func Assert(c testing.TB, obtained interface{}, checker Checker, args ...interface{}) {
	c.Helper()
	if !internalCheck(c, "Assert", obtained, checker, args...) {
		c.FailNow()
	}
}

func internalCheck(c testing.TB, funcName string, obtained interface{}, checker Checker, args ...interface{}) bool {
	c.Helper()
	if checker == nil {
		lines := []string{
			"",
			formatCaller(2),
			fmt.Sprintf("... %s(obtained, nil!?, ...):", funcName),
			"... Oops.. you've provided a nil checker!",
		}
		c.Error(strings.Join(lines, "\n"))
		return false
	}

	// If the last argument is a bug info, extract it out.
	var comment CommentInterface
	if len(args) > 0 {
		if c, ok := args[len(args)-1].(CommentInterface); ok {
			comment = c
			args = args[:len(args)-1]
		}
	}

	params := append([]interface{}{obtained}, args...)
	info := checker.Info()

	if len(params) != len(info.Params) {
		names := append([]string{info.Params[0], info.Name}, info.Params[1:]...)
		lines := []string{
			"",
			formatCaller(2),
			fmt.Sprintf("... %s(%s):", funcName, strings.Join(names, ", ")),
			fmt.Sprintf("... Wrong number of parameters for %s: want %d, got %d", info.Name, len(names), len(params)+1),
		}
		c.Error(strings.Join(lines, "\n"))
		return false
	}

	// Copy since it may be mutated by Check.
	names := append([]string{}, info.Params...)

	// Do the actual check.
	result, error := checker.Check(params, names)
	if !result || error != "" {
		lines := []string{
			"",
			formatCaller(2),
		}
		for i := 0; i != len(params); i++ {
			lines = append(lines, formatValue(names[i], params[i]))
		}
		if comment != nil {
			lines = append(lines, "... " + comment.CheckCommentString())
		}
		if error != "" {
			lines = append(lines, "... " + error)
		}
		c.Error(strings.Join(lines, "\n"))
		return false
	}
	return true
}

func formatCaller(skip int) string {
	// This is a bit heavier than it ought to be.
	skip++ // Our own frame.
	_, path, line, ok := runtime.Caller(skip)
	if !ok {
		return "    ..."
	}

	code, err := printLine(path, line)
	if code == "" {
		code = "..." // XXX Open the file and take the raw line.
		if err != nil {
			code += err.Error()
		}
	}
	return indent(code, "    ")
}

func formatValue(label string, value interface{}) string {
	if label == "" {
		if hasStringOrError(value) {
			return fmt.Sprintf("... %#v (%q)", value, value)
		} else {
			return fmt.Sprintf("... %#v", value)
		}
	} else if value == nil {
		return fmt.Sprintf("... %s = nil", label)
	} else {
		if hasStringOrError(value) {
			fv := fmt.Sprintf("%#v", value)
			qv := fmt.Sprintf("%q", value)
			if fv != qv {
				return fmt.Sprintf("... %s %s = %s (%s)", label, reflect.TypeOf(value), fv, qv)
			}
		}
		if s, ok := value.(string); ok && isMultiLine(s) {
			return fmt.Sprintf("... %s %s = \"\" +\n%s", label, reflect.TypeOf(value), formatMultiLine(s, true))
		} else {
			return fmt.Sprintf("... %s %s = %#v", label, reflect.TypeOf(value), value)
		}
	}
}

func hasStringOrError(x interface{}) (ok bool) {
	_, ok = x.(fmt.Stringer)
	if ok {
		return
	}
	_, ok = x.(error)
	return
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
	return b[:len(b)-1]
}

func isMultiLine(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '\n' {
			return true
		}
	}
	return false
}
