// These tests verify the inner workings of the helper methods associated
// with gocheck.T.

package gocheck_test

import (
    "gocheck"
    gocheck_local "gocheck/local"
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
    checkError string
    failCheck, noExpectedValue bool
    obtained, expected interface{}
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
// Ensure the internal interface matches the one in the subpackage.

// The Checker interface is internal inside the main gocheck package, to
// avoid importing the subpackage and thus having to know its location.
// Here we ensure that the interface in the subpackage actually matches,
// and thus that the whole thing works end to end.
func (s *HelpersS) TestCheckerInterface(c *gocheck.C) {
    testHelperSuccess(c, "Check(1, Equals, 1)", true, func() interface{} {
        return c.Check(1, gocheck_local.Equals, 1)
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
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Myexpected \\(int\\): 2\n\n"
    testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, 2)
    })
}

func (s *HelpersS) TestCheckFailWithExpectedAndBugInfo(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Myexpected \\(int\\): 2\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, 2, myBug("Hello world!"))
    })
}

func (s *HelpersS) TestCheckFailWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(myobtained, MyChecker\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n\n"
    testHelperFailure(c, "Check(1, checker)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker)
    })
}

func (s *HelpersS) TestCheckFailWithoutExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(myobtained, MyChecker\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Check(1, checker, msg)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, myBug("Hello world!"))
    })
}

func (s *HelpersS) TestCheckWithMissingExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
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
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
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
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Myexpected \\(int\\): 2\n" +
           "\\.+ Some not so cool data provided!\n\n"
    testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, 2)
    })
}

func (s *HelpersS) TestCheckWithNilChecker(c *gocheck.C) {
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
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
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Myexpected \\(int\\): 2\n\n"
    testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, 2)
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Myexpected \\(int\\): 2\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Assert(1, checker, 2, msg)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, 2, myBug("Hello world!"))
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(myobtained, MyChecker\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n\n"
    testHelperFailure(c, "Assert(1, checker)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker)
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithoutExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(myobtained, MyChecker\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Assert(1, checker, msg)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, myBug("Hello world!"))
        return nil
    })
}

func (s *HelpersS) TestAssertWithMissingExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
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
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(int\\): 1\n" +
           "\\.+ Myexpected \\(int\\): 2\n" +
           "\\.+ Some not so cool data provided!\n\n"
    testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, 2)
        return nil
    })
}

func (s *HelpersS) TestAssertWithNilChecker(c *gocheck.C) {
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
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
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ Check\\(myobtained, MyChecker, myexpected\\):\n" +
           "\\.+ Myobtained \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n" +
           "\\.+ Myexpected \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x3}\n\n"
    testHelperFailure(c, "Check([]byte{1}, chk, []byte{3})", false, false, log,
                      func() interface{} {
        return c.Check([]byte{1,2}, checker, []byte{1,3})
    })
}


// -----------------------------------------------------------------------
// Old tests for Assert*() and Check*() helpers not based on checkers.

func (s *HelpersS) TestCheckEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n\n"
    testHelperFailure(c, "CheckEqual(10, 20)", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, 20)
    })
}

func (s *HelpersS) TestCheckEqualFailingWithDiffTypes(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(uint\\): 0xa\n\n"
    testHelperFailure(c, "CheckEqual(10, uint(10))", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, uint(10))
    })
}

func (s *HelpersS) TestCheckEqualFailingWithNil(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(nil\\): nil\n\n"
    testHelperFailure(c, "CheckEqual(10, nil)", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, nil)
    })
}

func (s *HelpersS) TestCheckEqualWithBugInfo(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "CheckEqual(10, 20, issue)", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, 20, "That's clearly ", "WRONG!")
    })
}

func (s *HelpersS) TestCheckNotEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckNotEqual(10, 20)", true, func() interface{} {
        return c.CheckNotEqual(10, 20)
    })
}

func (s *HelpersS) TestCheckNotEqualSucceedingWithNil(c *gocheck.C) {
    testHelperSuccess(c, "CheckNotEqual(10, nil)", true, func() interface{} {
        return c.CheckNotEqual(10, nil)
    })
}

func (s *HelpersS) TestCheckNotEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n\n"
    testHelperFailure(c, "CheckNotEqual(10, 10)", false, false, log,
                      func() interface{} {
        return c.CheckNotEqual(10, 10)
    })
}

func (s *HelpersS) TestCheckNotEqualWithBugInfo(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "CheckNotEqual(10, 10, issue)", false, false, log,
                      func() interface{} {
        return c.CheckNotEqual(10, 10, "That's clearly ", "WRONG!")
    })
}

func (s *HelpersS) TestAssertEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertEqual(10, 10)", nil, func() interface{} {
        c.AssertEqual(10, 10)
        return nil
    })
}

func (s *HelpersS) TestAssertEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n\n"
    testHelperFailure(c, "AssertEqual(10, 20)", nil, true, log,
                      func() interface{} {
        c.AssertEqual(10, 20)
        return nil
    })
}

func (s *HelpersS) TestAssertEqualWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "AssertEqual(10, 20, issue)", nil, true, log,
                      func() interface{} {
        c.AssertEqual(10, 20, "That's clearly ", "WRONG!")
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertNotEqual(10, 20)", nil, func() interface{} {
        c.AssertNotEqual(10, 20)
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n\n"
    testHelperFailure(c, "AssertNotEqual(10, 10)", nil, true, log,
                      func() interface{} {
        c.AssertNotEqual(10, 10)
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "AssertNotEqual(10, 10, issue)", nil, true, log,
                      func() interface{} {
        c.AssertNotEqual(10, 10, "That's clearly ", "WRONG!")
        return nil
    })
}


func (s *HelpersS) TestCheckEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckEqual([]byte, []byte)", true, func() interface{} {
        return c.CheckEqual([]byte{1,2}, []byte{1,2})
    })
}

func (s *HelpersS) TestCheckEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n" +
           "\\.+ Expected \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x3}\n\n"
    testHelperFailure(c, "CheckEqual([]byte{2}, []byte{3})", false, false, log,
                      func() interface{} {
        return c.CheckEqual([]byte{1,2}, []byte{1,3})
    })
}

func (s *HelpersS) TestCheckNotEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckNotEqual([]byte, []byte)", true,
                      func() interface{} {
        return c.CheckNotEqual([]byte{1,2}, []byte{1,3})
    })
}

func (s *HelpersS) TestCheckNotEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n\n"
    testHelperFailure(c, "CheckNotEqual([]byte{2}, []byte{3})", false, false,
                      log, func() interface{} {
        return c.CheckNotEqual([]byte{1,2}, []byte{1,2})
    })
}


func (s *HelpersS) TestAssertEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertEqual([]byte, []byte)", nil,
                      func() interface{} {
        c.AssertEqual([]byte{1,2}, []byte{1,2})
        return nil
    })
}

func (s *HelpersS) TestAssertEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n" +
           "\\.+ Expected \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x3}\n\n"
    testHelperFailure(c, "AssertEqual([]byte{2}, []byte{3})", nil, true, log,
                      func() interface{} {
        c.AssertEqual([]byte{1,2}, []byte{1,3})
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertNotEqual([]byte, []byte)", nil,
                      func() interface{} {
        c.AssertNotEqual([]byte{1,2}, []byte{1,3})
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n\n"
    testHelperFailure(c, "AssertNotEqual([]byte{2}, []byte{3})", nil, true,
                      log, func() interface{} {
        c.AssertNotEqual([]byte{1,2}, []byte{1,2})
        return nil
    })
}

func (s *HelpersS) TestCheckMatchSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckErr('foo', 'fo*')", true, func() interface{} {
        return c.CheckMatch("foo", "fo*")
    })
}

func (s *HelpersS) TestCheckMatchFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foo\"\n" +
           "\\.+ Expected to match expression: \"bar\"\n\n"
    testHelperFailure(c, "CheckMatch('foo', 'bar')", false, false, log,
                      func() interface{} {
        return c.CheckMatch("foo", "bar")
    })
}

func (s *HelpersS) TestCheckMatchFailingWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foo\"\n" +
           "\\.+ Expected to match expression: \"bar\"\n" +
           "\\.+ Foo bar!\n\n"
    testHelperFailure(c, "CheckMatch('foo', 'bar')", false, false, log,
                      func() interface{} {
        return c.CheckMatch("foo", "bar", "Foo", " bar!")
    })
}


func (s *HelpersS) TestAssertMatchSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertMatch(s, exp)", nil, func() interface{} {
        c.AssertMatch("str error", "str.*r")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchSucceedingWithError(c *gocheck.C) {
    testHelperSuccess(c, "AssertMatch(os.Error, exp)", nil, func() interface{} {
        c.AssertMatch(os.Errno(13), "perm.*denied")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(os\\.Errno\\): 13 \\(permission denied\\)\n" +
           "\\.+ Expected to match expression: \"foo\"\n\n"
    testHelperFailure(c, "AssertMatch(error, foo)", nil, true, log,
                      func() interface{} {
        c.AssertMatch(os.Errno(13), "foo")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchFailingWithPureStrMatch(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foobar\"\n" +
           "\\.+ Expected to match expression: \"foobaz\"\n\n"
    testHelperFailure(c, "AssertMatch('foobar', 'foobaz')", nil, true, log,
                      func() interface{} {
        c.AssertMatch("foobar", "foobaz")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchFailingWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foobar\"\n" +
           "\\.+ Expected to match expression: \"foobaz\"\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "AssertMatch('foobar', 'foobaz')", nil, true, log,
                      func() interface{} {
        c.AssertMatch("foobar", "foobaz", "That's clearly ", "WRONG!")
        return nil
    })
}


// -----------------------------------------------------------------------
// MakeDir() tests.

type MkDirHelper struct {
    path1 string
    path2 string
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
    c.Assert(output.value, gocheck_local.Equals, "")
    c.Check(helper.isDir1, gocheck_local.Equals, true)
    c.Check(helper.isDir2, gocheck_local.Equals, true)
    c.Check(helper.isDir3, gocheck_local.Equals, true)
    c.Check(helper.isDir4, gocheck_local.Equals, true)
    c.Check(helper.path1, gocheck_local.Not(gocheck_local.Equals),
            helper.path2)
    c.Check(isDir(helper.path1), gocheck_local.Equals, false)
    c.Check(isDir(helper.path2), gocheck_local.Equals, false)
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
                        name: name,
                        result: expectedResult,
                        failed: false,
                        log: "",
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
                        name: name,
                        result: expectedResult,
                        failed: true,
                        log: log,
                   })
    })()
    result = closure()
    if shouldStop {
        c.Logf("%s didn't stop when it should", name)
    }
}
