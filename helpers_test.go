// These tests verify the inner workings of the helper methods associated
// with check.T.

package check

import (
	"os"
	"reflect"
	"runtime"
	"sync"
)

var helpersS = Suite(&HelpersS{})

type HelpersS struct{}

func (s *HelpersS) TestCountSuite(c *C) {
	suitesRun += 1
}

// -----------------------------------------------------------------------
// Fake checker and bug info to verify the behavior of Assert() and Check().

type MyChecker struct {
	info   *CheckerInfo
	params []interface{}
	names  []string
	result bool
	error  string
}

func (checker *MyChecker) Info() *CheckerInfo {
	if checker.info == nil {
		return &CheckerInfo{Name: "MyChecker", Params: []string{"myobtained", "myexpected"}}
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

func (s *HelpersS) TestCheckerInterface(c *C) {
	testHelperSuccess(c, "Check(1, Equals, 1)", true, func() interface{} {
		return c.Check(1, Equals, 1)
	})
}

// -----------------------------------------------------------------------
// Tests for Check(), mostly the same as for Assert() following these.

func (s *HelpersS) TestCheckSucceedWithExpected(c *C) {
	checker := &MyChecker{result: true}
	testHelperSuccess(c, "Check(1, checker, 2)", true, func() interface{} {
		return c.Check(1, checker, 2)
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1, 2}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestCheckSucceedWithoutExpected(c *C) {
	checker := &MyChecker{result: true, info: &CheckerInfo{Params: []string{"myvalue"}}}
	testHelperSuccess(c, "Check(1, checker)", true, func() interface{} {
		return c.Check(1, checker)
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestCheckFailWithExpected(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n\n"
	testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
		func() interface{} {
			return c.Check(1, checker, 2)
		})
}

func (s *HelpersS) TestCheckFailWithExpectedAndComment(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker, 2, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Hello world!\n\n"
	testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
		func() interface{} {
			return c.Check(1, checker, 2, myComment("Hello world!"))
		})
}

func (s *HelpersS) TestCheckFailWithExpectedAndStaticComment(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    // Nice leading comment\\.\n" +
		"    return c\\.Check\\(1, checker, 2\\) // Hello there\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n\n"
	testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
		func() interface{} {
			// Nice leading comment.
			return c.Check(1, checker, 2) // Hello there
		})
}

func (s *HelpersS) TestCheckFailWithoutExpected(c *C) {
	checker := &MyChecker{result: false, info: &CheckerInfo{Params: []string{"myvalue"}}}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker\\)\n" +
		"\\.+ myvalue int = 1\n\n"
	testHelperFailure(c, "Check(1, checker)", false, false, log,
		func() interface{} {
			return c.Check(1, checker)
		})
}

func (s *HelpersS) TestCheckFailWithoutExpectedAndMessage(c *C) {
	checker := &MyChecker{result: false, info: &CheckerInfo{Params: []string{"myvalue"}}}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myvalue int = 1\n" +
		"\\.+ Hello world!\n\n"
	testHelperFailure(c, "Check(1, checker, msg)", false, false, log,
		func() interface{} {
			return c.Check(1, checker, myComment("Hello world!"))
		})
}

func (s *HelpersS) TestCheckWithMissingExpected(c *C) {
	checker := &MyChecker{result: true}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker\\)\n" +
		"\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
		"\\.+ Wrong number of parameters for MyChecker: " +
		"want 3, got 2\n\n"
	testHelperFailure(c, "Check(1, checker, !?)", false, false, log,
		func() interface{} {
			return c.Check(1, checker)
		})
}

func (s *HelpersS) TestCheckWithTooManyExpected(c *C) {
	checker := &MyChecker{result: true}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker, 2, 3\\)\n" +
		"\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
		"\\.+ Wrong number of parameters for MyChecker: " +
		"want 3, got 4\n\n"
	testHelperFailure(c, "Check(1, checker, 2, 3)", false, false, log,
		func() interface{} {
			return c.Check(1, checker, 2, 3)
		})
}

func (s *HelpersS) TestCheckWithError(c *C) {
	checker := &MyChecker{result: false, error: "Some not so cool data provided!"}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Some not so cool data provided!\n\n"
	testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
		func() interface{} {
			return c.Check(1, checker, 2)
		})
}

func (s *HelpersS) TestCheckWithNilChecker(c *C) {
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, nil\\)\n" +
		"\\.+ Check\\(obtained, nil!\\?, \\.\\.\\.\\):\n" +
		"\\.+ Oops\\.\\. you've provided a nil checker!\n\n"
	testHelperFailure(c, "Check(obtained, nil)", false, false, log,
		func() interface{} {
			return c.Check(1, nil)
		})
}

func (s *HelpersS) TestCheckWithParamsAndNamesMutation(c *C) {
	checker := &MyChecker{result: false, params: []interface{}{3, 4}, names: []string{"newobtained", "newexpected"}}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    return c\\.Check\\(1, checker, 2\\)\n" +
		"\\.+ newobtained int = 3\n" +
		"\\.+ newexpected int = 4\n\n"
	testHelperFailure(c, "Check(1, checker, 2) with mutation", false, false, log,
		func() interface{} {
			return c.Check(1, checker, 2)
		})
}

// -----------------------------------------------------------------------
// Tests for Assert(), mostly the same as for Check() above.

func (s *HelpersS) TestAssertSucceedWithExpected(c *C) {
	checker := &MyChecker{result: true}
	testHelperSuccess(c, "Assert(1, checker, 2)", nil, func() interface{} {
		c.Assert(1, checker, 2)
		return nil
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1, 2}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestAssertSucceedWithoutExpected(c *C) {
	checker := &MyChecker{result: true, info: &CheckerInfo{Params: []string{"myvalue"}}}
	testHelperSuccess(c, "Assert(1, checker)", nil, func() interface{} {
		c.Assert(1, checker)
		return nil
	})
	if !reflect.DeepEqual(checker.params, []interface{}{1}) {
		c.Fatalf("Bad params for check: %#v", checker.params)
	}
}

func (s *HelpersS) TestAssertFailWithExpected(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n\n"
	testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
		func() interface{} {
			c.Assert(1, checker, 2)
			return nil
		})
}

func (s *HelpersS) TestAssertFailWithExpectedAndMessage(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, checker, 2, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Hello world!\n\n"
	testHelperFailure(c, "Assert(1, checker, 2, msg)", nil, true, log,
		func() interface{} {
			c.Assert(1, checker, 2, myComment("Hello world!"))
			return nil
		})
}

func (s *HelpersS) TestAssertFailWithoutExpected(c *C) {
	checker := &MyChecker{result: false, info: &CheckerInfo{Params: []string{"myvalue"}}}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, checker\\)\n" +
		"\\.+ myvalue int = 1\n\n"
	testHelperFailure(c, "Assert(1, checker)", nil, true, log,
		func() interface{} {
			c.Assert(1, checker)
			return nil
		})
}

func (s *HelpersS) TestAssertFailWithoutExpectedAndMessage(c *C) {
	checker := &MyChecker{result: false, info: &CheckerInfo{Params: []string{"myvalue"}}}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, checker, myComment\\(\"Hello world!\"\\)\\)\n" +
		"\\.+ myvalue int = 1\n" +
		"\\.+ Hello world!\n\n"
	testHelperFailure(c, "Assert(1, checker, msg)", nil, true, log,
		func() interface{} {
			c.Assert(1, checker, myComment("Hello world!"))
			return nil
		})
}

func (s *HelpersS) TestAssertWithMissingExpected(c *C) {
	checker := &MyChecker{result: true}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, checker\\)\n" +
		"\\.+ Assert\\(myobtained, MyChecker, myexpected\\):\n" +
		"\\.+ Wrong number of parameters for MyChecker: " +
		"want 3, got 2\n\n"
	testHelperFailure(c, "Assert(1, checker, !?)", nil, true, log,
		func() interface{} {
			c.Assert(1, checker)
			return nil
		})
}

func (s *HelpersS) TestAssertWithError(c *C) {
	checker := &MyChecker{result: false, error: "Some not so cool data provided!"}
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, checker, 2\\)\n" +
		"\\.+ myobtained int = 1\n" +
		"\\.+ myexpected int = 2\n" +
		"\\.+ Some not so cool data provided!\n\n"
	testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
		func() interface{} {
			c.Assert(1, checker, 2)
			return nil
		})
}

func (s *HelpersS) TestAssertWithNilChecker(c *C) {
	log := "(?s)helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
		"    c\\.Assert\\(1, nil\\)\n" +
		"\\.+ Assert\\(obtained, nil!\\?, \\.\\.\\.\\):\n" +
		"\\.+ Oops\\.\\. you've provided a nil checker!\n\n"
	testHelperFailure(c, "Assert(obtained, nil)", nil, true, log,
		func() interface{} {
			c.Assert(1, nil)
			return nil
		})
}

// -----------------------------------------------------------------------
// Ensure that values logged work properly in some interesting cases.

func (s *HelpersS) TestValueLoggingWithArrays(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test.go:[0-9]+:.*\nhelpers_test.go:[0-9]+:\n" +
		"    return c\\.Check\\(\\[\\]byte{1, 2}, checker, \\[\\]byte{1, 3}\\)\n" +
		"\\.+ myobtained \\[\\]uint8 = \\[\\]byte{0x1, 0x2}\n" +
		"\\.+ myexpected \\[\\]uint8 = \\[\\]byte{0x1, 0x3}\n\n"
	testHelperFailure(c, "Check([]byte{1}, chk, []byte{3})", false, false, log,
		func() interface{} {
			return c.Check([]byte{1, 2}, checker, []byte{1, 3})
		})
}

func (s *HelpersS) TestValueLoggingWithMultiLine(c *C) {
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test.go:[0-9]+:.*\nhelpers_test.go:[0-9]+:\n" +
		"    return c\\.Check\\(\"a\\\\nb\\\\n\", checker, \"a\\\\nb\\\\nc\"\\)\n" +
		"\\.+ myobtained string = \"\" \\+\n" +
		"\\.+     \"a\\\\n\" \\+\n" +
		"\\.+     \"b\\\\n\"\n" +
		"\\.+ myexpected string = \"\" \\+\n" +
		"\\.+     \"a\\\\n\" \\+\n" +
		"\\.+     \"b\\\\n\" \\+\n" +
		"\\.+     \"c\"\n\n"
	testHelperFailure(c, `Check("a\nb\n", chk, "a\nb\nc")`, false, false, log,
		func() interface{} {
			return c.Check("a\nb\n", checker, "a\nb\nc")
		})
}

func (s *HelpersS) TestValueLoggingWithMultiLineException(c *C) {
	// If the newline is at the end of the string, don't log as multi-line.
	checker := &MyChecker{result: false}
	log := "(?s)helpers_test.go:[0-9]+:.*\nhelpers_test.go:[0-9]+:\n" +
		"    return c\\.Check\\(\"a b\\\\n\", checker, \"a\\\\nb\"\\)\n" +
		"\\.+ myobtained string = \"a b\\\\n\"\n" +
		"\\.+ myexpected string = \"\" \\+\n" +
		"\\.+     \"a\\\\n\" \\+\n" +
		"\\.+     \"b\"\n\n"
	testHelperFailure(c, `Check("a b\n", chk, "a\nb")`, false, false, log,
		func() interface{} {
			return c.Check("a b\n", checker, "a\nb")
		})
}

// -----------------------------------------------------------------------
// MakeDir() tests.

type MkDirHelper struct {
	path1  string
	path2  string
	isDir1 bool
	isDir2 bool
	isDir3 bool
	isDir4 bool
}

func (s *MkDirHelper) SetUpSuite(c *C) {
	s.path1 = c.MkDir()
	s.isDir1 = isDir(s.path1)
}

func (s *MkDirHelper) Test(c *C) {
	s.path2 = c.MkDir()
	s.isDir2 = isDir(s.path2)
}

func (s *MkDirHelper) TearDownSuite(c *C) {
	s.isDir3 = isDir(s.path1)
	s.isDir4 = isDir(s.path2)
}

func (s *HelpersS) TestMkDir(c *C) {
	helper := MkDirHelper{}
	output := String{}
	Run(&helper, &RunConf{Output: &output})
	c.Assert(output.value, Equals, "")
	c.Check(helper.isDir1, Equals, true)
	c.Check(helper.isDir2, Equals, true)
	c.Check(helper.isDir3, Equals, true)
	c.Check(helper.isDir4, Equals, true)
	c.Check(helper.path1, Not(Equals),
		helper.path2)
	c.Check(isDir(helper.path1), Equals, false)
	c.Check(isDir(helper.path2), Equals, false)
}

func isDir(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

// Concurrent logging should not corrupt the underling buffer.
// Use go test -race to detect the race in this test.
func (s *HelpersS) TestConcurrentLogging(c *C) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(runtime.NumCPU()))
	var start, stop sync.WaitGroup
	start.Add(1)
	for i, n := 0, runtime.NumCPU()*2; i < n; i++ {
		stop.Add(1)
		go func(i int) {
			start.Wait()
			for j := 0; j < 30; j++ {
				c.Logf("Worker %d: line %d", i, j)
			}
			stop.Done()
		}(i)
	}
	start.Done()
	stop.Wait()
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

func (s *TestNameHelper) SetUpSuite(c *C)    { s.name1 = c.TestName() }
func (s *TestNameHelper) SetUpTest(c *C)     { s.name2 = c.TestName() }
func (s *TestNameHelper) Test(c *C)          { s.name3 = c.TestName() }
func (s *TestNameHelper) TearDownTest(c *C)  { s.name4 = c.TestName() }
func (s *TestNameHelper) TearDownSuite(c *C) { s.name5 = c.TestName() }

func (s *HelpersS) TestTestName(c *C) {
	helper := TestNameHelper{}
	output := String{}
	Run(&helper, &RunConf{Output: &output})
	c.Check(helper.name1, Equals, "")
	c.Check(helper.name2, Equals, "TestNameHelper.Test")
	c.Check(helper.name3, Equals, "TestNameHelper.Test")
	c.Check(helper.name4, Equals, "TestNameHelper.Test")
	c.Check(helper.name5, Equals, "")
}

// -----------------------------------------------------------------------
// A couple of helper functions to test helper functions. :-)

func testHelperSuccess(c *C, name string, expectedResult interface{}, closure func() interface{}) {
	var result interface{}
	defer (func() {
		if err := recover(); err != nil {
			panic(err)
		}
		checkState(c, result,
			&expectedState{
				name:   name,
				result: expectedResult,
				failed: false,
				log:    "",
			})
	})()
	result = closure()
}

func testHelperFailure(c *C, name string, expectedResult interface{}, shouldStop bool, log string, closure func() interface{}) {
	var result interface{}
	defer (func() {
		if err := recover(); err != nil {
			panic(err)
		}
		checkState(c, result,
			&expectedState{
				name:   name,
				result: expectedResult,
				failed: true,
				log:    log,
			})
	})()
	result = closure()
	if shouldStop {
		c.Logf("%s didn't stop when it should", name)
	}
}
