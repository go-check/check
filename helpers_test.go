// These tests verify the inner workings of the helper methods associated
// with gocheck.T.

package gocheck_test

import (
    "gocheck"
    "os"
)


var helpersS = gocheck.Suite(&HelpersS{})

type HelpersS struct{}

func (s *HelpersS) TestCountSuite(c *gocheck.C) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Fake checker and bug info to verify the behavior of Assert() and Check().

type MyChecker struct {
    checkError                 string
    failCheck, noExpectedValue bool
    obtained, expected         interface{}
}

func (checker *MyChecker) Name() string {
    return "MyChecker"
}

func (checker *MyChecker) VarNames() (obtained, expected string) {
    return "myobtained", "myexpected"
}

func (checker *MyChecker) NeedsExpectedValue() bool {
    return !checker.noExpectedValue
}

func (checker *MyChecker) Check(obtained, expected interface{}) (bool, string) {
    checker.obtained = obtained
    checker.expected = expected
    return !checker.failCheck, checker.checkError
}


type myBugInfo struct {
    info string
}

func (bug *myBugInfo) GetBugInfo() string {
    return bug.info
}

func myBug(info string) *myBugInfo {
    return &myBugInfo{info}
}


// -----------------------------------------------------------------------
// Ensure a real checker actually works fine.

func (s *HelpersS) TestCheckerInterface(c *gocheck.C) {
    testHelperSuccess(c, "Check(1, Equals, 1)", true, func() interface{} {
        return c.Check(1, gocheck.Equals, 1)
    })
}


// -----------------------------------------------------------------------
// Tests for Check(), mostly the same as for Assert() following these.

func (s *HelpersS) TestCheckSucceedWithExpected(c *gocheck.C) {
    checker := &MyChecker{}
    testHelperSuccess(c, "Check(1, checker, 2)", true, func() interface{} {
        return c.Check(1, checker, 2)
    })
    if checker.obtained != 1 || checker.expected != 2 {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
            checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestCheckSucceedWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{noExpectedValue: true}
    testHelperSuccess(c, "Check(1, checker)", true, func() interface{} {
        return c.Check(1, checker)
    })
    if checker.obtained != 1 || checker.expected != nil {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
            checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestCheckFailWithExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker, 2\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ myexpected int = 2\n\n"
    testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
        func() interface{} {
            return c.Check(1, checker, 2)
        })
}

func (s *HelpersS) TestCheckFailWithExpectedAndBugInfo(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker, 2, myBug\\(\"Hello world!\"\\)\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ myexpected int = 2\n" +
        "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
        func() interface{} {
            return c.Check(1, checker, 2, myBug("Hello world!"))
        })
}

func (s *HelpersS) TestCheckFailWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker\\)\n" +
        "\\.+ myobtained int = 1\n\n"
    testHelperFailure(c, "Check(1, checker)", false, false, log,
        func() interface{} {
            return c.Check(1, checker)
        })
}

func (s *HelpersS) TestCheckFailWithoutExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker, myBug\\(\"Hello world!\"\\)\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Check(1, checker, msg)", false, false, log,
        func() interface{} {
            return c.Check(1, checker, myBug("Hello world!"))
        })
}

func (s *HelpersS) TestCheckWithMissingExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker\\)\n" +
        "\\.+ Check\\(myobtained, MyChecker, >myexpected<\\):\n" +
        "\\.+ Wrong number of myexpected args for MyChecker: " +
        "want 1, got 0\n\n"
    testHelperFailure(c, "Check(1, checker, !?)", false, false, log,
        func() interface{} {
            return c.Check(1, checker)
        })
}

func (s *HelpersS) TestCheckWithTooManyExpected(c *gocheck.C) {
    checker := &MyChecker{noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker, 1\\)\n" +
        "\\.+ Check\\(myobtained, MyChecker, >myexpected<\\):\n" +
        "\\.+ Wrong number of myexpected args for MyChecker: " +
        "want 0, got 1\n\n"
    testHelperFailure(c, "Check(1, checker, !?)", false, false, log,
        func() interface{} {
            return c.Check(1, checker, 1)
        })
}

func (s *HelpersS) TestCheckWithError(c *gocheck.C) {
    checker := &MyChecker{checkError: "Some not so cool data provided!"}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, checker, 2\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ myexpected int = 2\n" +
        "\\.+ Some not so cool data provided!\n\n"
    testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
        func() interface{} {
            return c.Check(1, checker, 2)
        })
}

func (s *HelpersS) TestCheckWithNilChecker(c *gocheck.C) {
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    return c\\.Check\\(1, nil\\)\n" +
        "\\.+ Check\\(obtained, nil!\\?, \\.\\.\\.\\):\n" +
        "\\.+ Oops\\.\\. you've provided a nil checker!\n\n"
    testHelperFailure(c, "Check(obtained, nil)", false, false, log,
        func() interface{} {
            return c.Check(1, nil)
        })
}


// -----------------------------------------------------------------------
// Tests for Assert(), mostly the same as for Check() above.

func (s *HelpersS) TestAssertSucceedWithExpected(c *gocheck.C) {
    checker := &MyChecker{}
    testHelperSuccess(c, "Assert(1, checker, 2)", nil, func() interface{} {
        c.Assert(1, checker, 2)
        return nil
    })
    if checker.obtained != 1 || checker.expected != 2 {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
            checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestAssertSucceedWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{noExpectedValue: true}
    testHelperSuccess(c, "Assert(1, checker)", nil, func() interface{} {
        c.Assert(1, checker)
        return nil
    })
    if checker.obtained != 1 || checker.expected != nil {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
            checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestAssertFailWithExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    c\\.Assert\\(1, checker, 2\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ myexpected int = 2\n\n"
    testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
        func() interface{} {
            c.Assert(1, checker, 2)
            return nil
        })
}

func (s *HelpersS) TestAssertFailWithExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    c\\.Assert\\(1, checker, 2, myBug\\(\"Hello world!\"\\)\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ myexpected int = 2\n" +
        "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Assert(1, checker, 2, msg)", nil, true, log,
        func() interface{} {
            c.Assert(1, checker, 2, myBug("Hello world!"))
            return nil
        })
}

func (s *HelpersS) TestAssertFailWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    c\\.Assert\\(1, checker\\)\n" +
        "\\.+ myobtained int = 1\n\n"
    testHelperFailure(c, "Assert(1, checker)", nil, true, log,
        func() interface{} {
            c.Assert(1, checker)
            return nil
        })
}

func (s *HelpersS) TestAssertFailWithoutExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    c\\.Assert\\(1, checker, myBug\\(\"Hello world!\"\\)\\)\n" +
        "\\.+ myobtained int = 1\n" +
        "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Assert(1, checker, msg)", nil, true, log,
        func() interface{} {
            c.Assert(1, checker, myBug("Hello world!"))
            return nil
        })
}

func (s *HelpersS) TestAssertWithMissingExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
        "    c\\.Assert\\(1, checker\\)\n" +
        "\\.+ Assert\\(myobtained, MyChecker, >myexpected<\\):\n" +
        "\\.+ Wrong number of myexpected args for MyChecker: " +
        "want 1, got 0\n\n"
    testHelperFailure(c, "Assert(1, checker, !?)", nil, true, log,
        func() interface{} {
            c.Assert(1, checker)
            return nil
        })
}

func (s *HelpersS) TestAssertWithError(c *gocheck.C) {
    checker := &MyChecker{checkError: "Some not so cool data provided!"}
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
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

func (s *HelpersS) TestAssertWithNilChecker(c *gocheck.C) {
    log := "helpers_test\\.go:[0-9]+:.*\nhelpers_test\\.go:[0-9]+:\n" +
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

func (s *HelpersS) TestValueLoggingWithArrays(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test.go:[0-9]+:.*\nhelpers_test.go:[0-9]+:\n" +
        "    return c\\.Check\\(\\[\\]byte{1, 2}, checker, \\[\\]byte{1, 3}\\)\n" +
        "\\.+ myobtained \\[\\]uint8 = \\[\\]byte{0x1, 0x2}\n" +
        "\\.+ myexpected \\[\\]uint8 = \\[\\]byte{0x1, 0x3}\n\n"
    testHelperFailure(c, "Check([]byte{1}, chk, []byte{3})", false, false, log,
        func() interface{} {
            return c.Check([]byte{1, 2}, checker, []byte{1, 3})
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

func (s *MkDirHelper) SetUpSuite(c *gocheck.C) {
    s.path1 = c.MkDir()
    s.isDir1 = isDir(s.path1)
}

func (s *MkDirHelper) Test(c *gocheck.C) {
    s.path2 = c.MkDir()
    s.isDir2 = isDir(s.path2)
}

func (s *MkDirHelper) TearDownSuite(c *gocheck.C) {
    s.isDir3 = isDir(s.path1)
    s.isDir4 = isDir(s.path2)
}


func (s *HelpersS) TestMkDir(c *gocheck.C) {
    helper := MkDirHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Assert(output.value, gocheck.Equals, "")
    c.Check(helper.isDir1, gocheck.Equals, true)
    c.Check(helper.isDir2, gocheck.Equals, true)
    c.Check(helper.isDir3, gocheck.Equals, true)
    c.Check(helper.isDir4, gocheck.Equals, true)
    c.Check(helper.path1, gocheck.Not(gocheck.Equals),
        helper.path2)
    c.Check(isDir(helper.path1), gocheck.Equals, false)
    c.Check(isDir(helper.path2), gocheck.Equals, false)
}

func isDir(path string) bool {
    if stat, err := os.Stat(path); err == nil {
        return stat.IsDirectory()
    }
    return false
}


// -----------------------------------------------------------------------
// A couple of helper functions to test helper functions. :-)

func testHelperSuccess(c *gocheck.C, name string,
expectedResult interface{},
closure func() interface{}) {
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

func testHelperFailure(c *gocheck.C, name string,
expectedResult interface{},
shouldStop bool, log string,
closure func() interface{}) {
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
