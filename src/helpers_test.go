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

func (s *HelpersS) TestCheckEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckEqual(10, 10)", true, func() interface{} {
        return c.CheckEqual(10, 10)
    })
}

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

func (s *HelpersS) TestCheckEqualWithMessage(c *gocheck.C) {
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

func (s *HelpersS) TestCheckNotEqualWithMessage(c *gocheck.C) {
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
    c.AssertEqual(output.value, "")
    c.CheckEqual(helper.isDir1, true)
    c.CheckEqual(helper.isDir2, true)
    c.CheckEqual(helper.isDir3, true)
    c.CheckEqual(helper.isDir4, true)
    c.CheckNotEqual(helper.path1, helper.path2)
    c.CheckEqual(isDir(helper.path1), false)
    c.CheckEqual(isDir(helper.path2), false)
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
