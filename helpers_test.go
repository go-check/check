// These tests verify the inner workings of the helper methods associated
// with check.T.

package check_test

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	"gopkg.in/check.v1"
)

var helpersS = check.Suite(&HelpersS{})

type HelpersS struct{}

func (s *HelpersS) TestCountSuite(c *check.C) {
	suitesRun += 1
}

// -----------------------------------------------------------------------
// Fake checker and bug info to verify the behavior of Assert() and Check().

type MyChecker struct {
	info   *check.CheckerInfo
	params []interface{}
	names  []string
	result bool
	error  string
}

func (checker *MyChecker) Info() *check.CheckerInfo {
	if checker.info == nil {
		return &check.CheckerInfo{Name: "MyChecker", Params: []string{"myobtained", "myexpected"}}
	}
	return checker.info
}

func (checker *MyChecker) Check(params []interface{}, names []string) (bool, string) {
	rparams := checker.params
	rnames := checker.names
	checker.params = append([]interface{}{}, params...)
	checker.names = append([]string{}, names...)
	if rparams != nil {
		copy(params, rparams)
	}
	if rnames != nil {
		copy(names, rnames)
	}
	return checker.result, checker.error
}

type myCommentType string

func (c myCommentType) CheckCommentString() string {
	return string(c)
}

func myComment(s string) myCommentType {
	return myCommentType(s)
}

// -----------------------------------------------------------------------
// Ensure a real checker actually works fine.

func (s *HelpersS) TestCheckerInterface(c *check.C) {
	testHelperSuccess(c, "Check(1, Equals, 1)", true, func(t testing.TB) interface{} {
		return check.Check(t, 1, check.Equals, 1)
	})
}

// -----------------------------------------------------------------------
// Tests for Check(), mostly the same as for Assert() following these.

func (s *HelpersS) TestCheckSucceedWithExpected(c *check.C) {
	checker := &MyChecker{result: true}
	testHelperSuccess(c, "Check(1, checker, 2)", true, func(t testing.TB) interface{} {
		return check.Check(t, 1, checker, 2)
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1, 2}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestCheckSucceedWithoutExpected(c *check.C) {
	checker := &MyChecker{result: true, info: &check.CheckerInfo{Params: []string{"myvalue"}}}
	testHelperSuccess(c, "Check(1, checker)", true, func(t testing.TB) interface{} {
		return check.Check(t, 1, checker)
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestCheckFailWithExpected(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n"
	testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker, 2)
		})
}

func (s *HelpersS) TestCheckFailWithExpectedAndComment(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker, 2, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Hello world!\n"
	testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker, 2, myComment("Hello world!"))
		})
}

func (s *HelpersS) TestCheckFailWithExpectedAndStaticComment(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    // Nice leading comment\\.\n" +
		"    return check\\.Check\\(t, 1, checker, 2\\) // Hello there\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n"
	testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
		func(t testing.TB) interface{} {
			// Nice leading comment.
			return check.Check(t, 1, checker, 2) // Hello there
		})
}

func (s *HelpersS) TestCheckFailWithoutExpected(c *check.C) {
	checker := &MyChecker{result: false, info: &check.CheckerInfo{Params: []string{"myvalue"}}}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker\\)\n" +
		"\\.+ myvalue int = 1\n"
	testHelperFailure(c, "Check(1, checker)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker)
		})
}

func (s *HelpersS) TestCheckFailWithoutExpectedAndMessage(c *check.C) {
	checker := &MyChecker{result: false, info: &check.CheckerInfo{Params: []string{"myvalue"}}}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myvalue int = 1\n" +
		"\\.+ Hello world!\n"
	testHelperFailure(c, "Check(1, checker, msg)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker, myComment("Hello world!"))
		})
}

func (s *HelpersS) TestCheckWithMissingExpected(c *check.C) {
	checker := &MyChecker{result: true}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker\\)\n" +
		"\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
		"\\.+ Wrong number of parameters for MyChecker: " +
		"want 3, got 2\n"
	testHelperFailure(c, "Check(1, checker, !?)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker)
		})
}

func (s *HelpersS) TestCheckWithTooManyExpected(c *check.C) {
	checker := &MyChecker{result: true}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker, 2, 3\\)\n" +
		"\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
		"\\.+ Wrong number of parameters for MyChecker: " +
		"want 3, got 4\n"
	testHelperFailure(c, "Check(1, checker, 2, 3)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker, 2, 3)
		})
}

func (s *HelpersS) TestCheckWithError(c *check.C) {
	checker := &MyChecker{result: false, error: "Some not so cool data provided!"}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Some not so cool data provided!\n"
	testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker, 2)
		})
}

func (s *HelpersS) TestCheckWithNilChecker(c *check.C) {
	log := "\n" +
		"    return check\\.Check\\(t, 1, nil\\)\n" +
		"\\.+ Check\\(obtained, nil!\\?, \\.\\.\\.\\):\n" +
		"\\.+ Oops\\.\\. you've provided a nil checker!\n"
	testHelperFailure(c, "Check(obtained, nil)", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, nil)
		})
}

func (s *HelpersS) TestCheckWithParamsAndNamesMutation(c *check.C) {
	checker := &MyChecker{result: false, params: []interface{}{3, 4}, names: []string{"newobtained", "newexpected"}}
	log := "\n" +
		"    return check\\.Check\\(t, 1, checker, 2\\)\n" +
		"\\.+ newobtained int = 3\n" +
		"\\.+ newexpected int = 4\n"
	testHelperFailure(c, "Check(1, checker, 2) with mutation", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, 1, checker, 2)
		})
}

// -----------------------------------------------------------------------
// Tests for Assert(), mostly the same as for Check() above.

func (s *HelpersS) TestAssertSucceedWithExpected(c *check.C) {
	checker := &MyChecker{result: true}
	testHelperSuccess(c, "Assert(1, checker, 2)", nil, func(t testing.TB) interface{} {
		check.Assert(t, 1, checker, 2)
		return nil
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1, 2}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestAssertSucceedWithoutExpected(c *check.C) {
	checker := &MyChecker{result: true, info: &check.CheckerInfo{Params: []string{"myvalue"}}}
	testHelperSuccess(c, "Assert(1, checker)", nil, func(t testing.TB) interface{} {
		check.Assert(t, 1, checker)
		return nil
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestAssertFailWithExpected(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    check\\.Assert\\(t, 1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n"
	testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, checker, 2)
			return nil
		})
}

func (s *HelpersS) TestAssertFailWithExpectedAndMessage(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    check\\.Assert\\(t, 1, checker, 2, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Hello world!\n"
	testHelperFailure(c, "Assert(1, checker, 2, msg)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, checker, 2, myComment("Hello world!"))
			return nil
		})
}

func (s *HelpersS) TestAssertFailWithoutExpected(c *check.C) {
	checker := &MyChecker{result: false, info: &check.CheckerInfo{Params: []string{"myvalue"}}}
	log := "\n" +
		"    check\\.Assert\\(t, 1, checker\\)\n" +
		"\\.+ myvalue int = 1\n"
	testHelperFailure(c, "Assert(1, checker)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, checker)
			return nil
		})
}

func (s *HelpersS) TestAssertFailWithoutExpectedAndMessage(c *check.C) {
	checker := &MyChecker{result: false, info: &check.CheckerInfo{Params: []string{"myvalue"}}}
	log := "\n" +
		"    check\\.Assert\\(t, 1, checker, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myvalue int = 1\n" +
		"\\.+ Hello world!\n"
	testHelperFailure(c, "Assert(1, checker, msg)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, checker, myComment("Hello world!"))
			return nil
		})
}

func (s *HelpersS) TestAssertWithMissingExpected(c *check.C) {
	checker := &MyChecker{result: true}
	log := "\n" +
		"    check\\.Assert\\(t, 1, checker\\)\n" +
		"\\.+ Assert\\(myobtained, MyChecker, myexpected\\):\n" +
		"\\.+ Wrong number of parameters for MyChecker: " +
		"want 3, got 2\n"
	testHelperFailure(c, "Assert(1, checker, !?)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, checker)
			return nil
		})
}

func (s *HelpersS) TestAssertWithError(c *check.C) {
	checker := &MyChecker{result: false, error: "Some not so cool data provided!"}
	log := "\n" +
		"    check\\.Assert\\(t, 1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Some not so cool data provided!\n"
	testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, checker, 2)
			return nil
		})
}

func (s *HelpersS) TestAssertWithNilChecker(c *check.C) {
	log := "\n" +
		"    check\\.Assert\\(t, 1, nil\\)\n" +
		"\\.+ Assert\\(obtained, nil!\\?, \\.\\.\\.\\):\n" +
		"\\.+ Oops\\.\\. you've provided a nil checker!\n"
	testHelperFailure(c, "Assert(obtained, nil)", nil, true, log,
		func(t testing.TB) interface{} {
			check.Assert(t, 1, nil)
			return nil
		})
}

// -----------------------------------------------------------------------
// Ensure that values logged work properly in some interesting cases.

func (s *HelpersS) TestValueLoggingWithArrays(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    return check\\.Check\\(t, \\[\\]byte{1, 2}, checker, \\[\\]byte{1, 3}\\)\n" +
		"\\.+ myobtained \\[\\]uint8 = \\[\\]byte{0x1, 0x2}\n" +
		"\\.+ myexpected \\[\\]uint8 = \\[\\]byte{0x1, 0x3}\n"
	testHelperFailure(c, "Check([]byte{1}, chk, []byte{3})", false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, []byte{1, 2}, checker, []byte{1, 3})
		})
}

func (s *HelpersS) TestValueLoggingWithMultiLine(c *check.C) {
	checker := &MyChecker{result: false}
	log := "\n" +
		"    return check\\.Check\\(t, \"a\\\\nb\\\\n\", checker, \"a\\\\nb\\\\nc\"\\)\n" +
		"\\.+ myobtained string = \"\" \\+\n" +
		"\\.+     \"a\\\\n\" \\+\n" +
		"\\.+     \"b\\\\n\"\n" +
		"\\.+ myexpected string = \"\" \\+\n" +
		"\\.+     \"a\\\\n\" \\+\n" +
		"\\.+     \"b\\\\n\" \\+\n" +
		"\\.+     \"c\"\n"
	testHelperFailure(c, `Check("a\nb\n", chk, "a\nb\nc")`, false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, "a\nb\n", checker, "a\nb\nc")
		})
}

func (s *HelpersS) TestValueLoggingWithMultiLineException(c *check.C) {
	// If the newline is at the end of the string, don't log as multi-line.
	checker := &MyChecker{result: false}
	log := "\n" +
		"    return check\\.Check\\(t, \"a b\\\\n\", checker, \"a\\\\nb\"\\)\n" +
		"\\.+ myobtained string = \"a b\\\\n\"\n" +
		"\\.+ myexpected string = \"\" \\+\n" +
		"\\.+     \"a\\\\n\" \\+\n" +
		"\\.+     \"b\"\n"
	testHelperFailure(c, `Check("a b\n", chk, "a\nb")`, false, false, log,
		func(t testing.TB) interface{} {
			return check.Check(t, "a b\n", checker, "a\nb")
		})
}

// -----------------------------------------------------------------------
// Test the TestName function

type TestNameHelper struct {
	name1 string
	name2 string
	name3 string
	name4 string
	name5 string
}

func (s *TestNameHelper) SetUpSuite(c *check.C)    { s.name1 = c.TestName() }
func (s *TestNameHelper) SetUpTest(c *check.C)     { s.name2 = c.TestName() }
func (s *TestNameHelper) Test(c *check.C)          { s.name3 = c.TestName() }
func (s *TestNameHelper) TearDownTest(c *check.C)  { s.name4 = c.TestName() }
func (s *TestNameHelper) TearDownSuite(c *check.C) { s.name5 = c.TestName() }

func (s *HelpersS) TestTestName(c *check.C) {
	helper := TestNameHelper{}
	base := c.TestName()
	check.Run(c.T, &helper)
	c.Check(helper.name1, check.Equals, base+"/TestNameHelper")
	c.Check(helper.name2, check.Equals, base+"/TestNameHelper/Test")
	c.Check(helper.name3, check.Equals, base+"/TestNameHelper/Test")
	c.Check(helper.name4, check.Equals, base+"/TestNameHelper/Test")
	c.Check(helper.name5, check.Equals, base+"/TestNameHelper")
}

// -----------------------------------------------------------------------
// A couple of helper functions to test helper functions. :-)

// fakeT is a testing.TB that provides enough functionality for the
// Check and Asseert helpers to run.
//
// We embed a testing.TB value so that fakeT also implements the
// interface (despite the private method). Any methods we don't
// implement will pass through to the embedded testing.TB and panic,
// since the value is nil.
type fakeT struct {
	testing.TB

	log       bytes.Buffer
	result    interface{}
	failed    bool
	completed bool
}

func (t *fakeT) Helper() {}

func (t *fakeT) Error(args ...interface{}) {
	fmt.Fprintln(&t.log, args...)
	t.failed = true
}

func (t *fakeT) FailNow() {
	t.failed = true
	runtime.Goexit()
}

func (t *fakeT) run(closure func(t testing.TB) interface{}) {
	// run the closure in a goroutine so we can use Goexit to
	// implement FailNow.
	finish := make(chan struct{})
	go func() {
		defer close(finish)
		t.result = closure(t)
		t.completed = true
	}()
	<-finish
}

type expectedState struct {
	name      string
	result    interface{}
	failed    bool
	completed bool
	log       string
}

// Helper which checks the state of the test and ensures that it matches
// the given expectations.  Depends on c.Errorf() working, so shouldn't
// be used to test this one function.
func (t *fakeT) checkState(c *check.C, expected *expectedState) {
	c.Helper()
	log := t.log.String()
	matched, matchError := regexp.MatchString("^"+expected.log+"$", log)
	if matchError != nil {
		c.Errorf("Error in matching expression used in testing %s: %v",
			expected.name, matchError)
	} else if !matched {
		c.Errorf("%s logged:\n----------\n%s----------\n\nExpected:\n----------\n%s\n----------",
			expected.name, log, expected.log)
	}
	if t.result != expected.result {
		c.Errorf("%s returned %#v rather than %#v",
			expected.name, t.result, expected.result)
	}
	if t.failed != expected.failed {
		if t.failed {
			c.Errorf("%s has failed when it shouldn't", expected.name)
		} else {
			c.Errorf("%s has not failed when it should", expected.name)
		}
	}
	if t.completed != expected.completed {
		if t.completed {
			c.Errorf("%s has completed when it shouldn't", expected.name)
		} else {
			c.Errorf("%s has not completed when it should", expected.name)
		}
	}
}

func testHelperSuccess(c *check.C, name string, expectedResult interface{}, closure func(t testing.TB) interface{}) {
	c.Helper()
	// Run the closure in a goroutine with its a fake testing value
	var t fakeT
	t.run(closure)
	t.checkState(c, &expectedState{
		name:      name,
		result:    expectedResult,
		failed:    false,
		completed: true,
		log:       "",
	})
}

func testHelperFailure(c *check.C, name string, expectedResult interface{}, shouldStop bool, log string, closure func(t testing.TB) interface{}) {
	c.Helper()
	// Run the closure in a goroutine with its a fake testing value
	var t fakeT
	t.run(closure)
	t.checkState(c, &expectedState{
		name:      name,
		result:    expectedResult,
		failed:    true,
		completed: !shouldStop,
		log:       log,
	})
}
