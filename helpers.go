package check

import (
	"fmt"
	"strings"
	"time"
)

// TestName returns the current test name in the form "SuiteName.TestName"
func (c *C) TestName() string {
	return c.Name()
}

// -----------------------------------------------------------------------
// Basic succeeding/failing logic.

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
	return c.internalCheck("Check", obtained, checker, args...)
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
	if !c.internalCheck("Assert", obtained, checker, args...) {
		c.FailNow()
	}
}

func (c *C) internalCheck(funcName string, obtained interface{}, checker Checker, args ...interface{}) bool {
	c.Helper()
	if checker == nil {
		c.logCaller(2)
		c.logString(fmt.Sprintf("%s(obtained, nil!?, ...):", funcName))
		c.logString("Oops.. you've provided a nil checker!")
		c.logNewLine()
		c.Fail()
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
		c.logCaller(2)
		c.logString(fmt.Sprintf("%s(%s):", funcName, strings.Join(names, ", ")))
		c.logString(fmt.Sprintf("Wrong number of parameters for %s: want %d, got %d", info.Name, len(names), len(params)+1))
		c.logNewLine()
		c.Fail()
		return false
	}

	// Copy since it may be mutated by Check.
	names := append([]string{}, info.Params...)

	// Do the actual check.
	result, error := checker.Check(params, names)
	if !result || error != "" {
		c.logCaller(2)
		for i := 0; i != len(params); i++ {
			c.logValue(names[i], params[i])
		}
		if comment != nil {
			c.logString(comment.CheckCommentString())
		}
		if error != "" {
			c.logString(error)
		}
		c.logNewLine()
		c.Fail()
		return false
	}
	return true
}
