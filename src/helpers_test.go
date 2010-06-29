// These tests verify the inner workings of the helper methods associated
// with gocheck.T.

package gocheck_test

import (
    "gocheck"
    "testing"
)


func TestHelpers(t *testing.T) {
    gocheck.RunTestingT(&HelpersS{}, t)
}


// -----------------------------------------------------------------------
// Helpers test suite.

type HelpersS struct{}

func (s *HelpersS) TestCheckEqualSucceeding(t *gocheck.T) {
    testHelperSuccess(t, "CheckEqual(10, 10)", true, func() interface{} {
        return t.CheckEqual(10, 10)
    })
}

func (s *HelpersS) TestCheckEqualFailing(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:CheckEqual\\(A, B\\): A != B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 20\n\n"
    testHelperFailure(t, "CheckEqual(10, 20)", false, false, log,
                      func() interface{} {
        return t.CheckEqual(10, 20)
    })
}

func (s *HelpersS) TestCheckEqualWithMessage(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:CheckEqual\\(A, B\\): A != B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 20\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(t, "CheckEqual(10, 20, issue)", false, false, log,
                      func() interface{} {
        return t.CheckEqual(10, 20, "That's clearly ", "WRONG!")
    })
}

func (s *HelpersS) TestCheckNotEqualSucceeding(t *gocheck.T) {
    testHelperSuccess(t, "CheckNotEqual(10, 20)", true, func() interface{} {
        return t.CheckNotEqual(10, 20)
    })
}

func (s *HelpersS) TestCheckNotEqualFailing(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:CheckNotEqual\\(A, B\\): A == B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 10\n\n"
    testHelperFailure(t, "CheckNotEqual(10, 10)", false, false, log,
                      func() interface{} {
        return t.CheckNotEqual(10, 10)
    })
}

func (s *HelpersS) TestCheckNotEqualWithMessage(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:CheckNotEqual\\(A, B\\): A == B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 10\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(t, "CheckNotEqual(10, 10, issue)", false, false, log,
                      func() interface{} {
        return t.CheckNotEqual(10, 10, "That's clearly ", "WRONG!")
    })
}

func (s *HelpersS) TestAssertEqualSucceeding(t *gocheck.T) {
    testHelperSuccess(t, "AssertEqual(10, 10)", nil, func() interface{} {
        t.AssertEqual(10, 10)
        return nil
    })
}

func (s *HelpersS) TestAssertEqualFailing(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:AssertEqual\\(A, B\\): A != B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 20\n\n"
    testHelperFailure(t, "AssertEqual(10, 20)", nil, true, log,
                      func() interface{} {
        t.AssertEqual(10, 20)
        return nil
    })
}

func (s *HelpersS) TestAssertEqualWithMessage(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:AssertEqual\\(A, B\\): A != B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 20\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(t, "AssertEqual(10, 20, issue)", nil, true, log,
                      func() interface{} {
        t.AssertEqual(10, 20, "That's clearly ", "WRONG!")
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualSucceeding(t *gocheck.T) {
    testHelperSuccess(t, "AssertNotEqual(10, 20)", nil, func() interface{} {
        t.AssertNotEqual(10, 20)
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualFailing(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:AssertNotEqual\\(A, B\\): A == B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 10\n\n"
    testHelperFailure(t, "AssertNotEqual(10, 10)", nil, true, log,
                      func() interface{} {
        t.AssertNotEqual(10, 10)
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualWithMessage(t *gocheck.T) {
    log := "\n\\.+ [0-9]+:AssertNotEqual\\(A, B\\): A == B\n" +
           "\\.+ A: 10\n" +
           "\\.+ B: 10\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(t, "AssertNotEqual(10, 10, issue)", nil, true, log,
                      func() interface{} {
        t.AssertNotEqual(10, 10, "That's clearly ", "WRONG!")
        return nil
    })
}


// -----------------------------------------------------------------------
// A couple of helper functions to test helper functions. :-)

func testHelperSuccess(t *gocheck.T, name string,
                       expectedResult interface{},
                       closure func() interface{}) {
    var result interface{}
    defer (func() {
        if err := recover(); err != nil {
            panic(err)
        }
        checkState(t, result,
                   &expectedState{
                        name: name,
                        result: expectedResult,
                        failed: false,
                        log: "",
                   })
    })()
    result = closure()
}

func testHelperFailure(t *gocheck.T, name string,
                       expectedResult interface{},
                       shouldStop bool, log string,
                       closure func() interface{}) {
    var result interface{}
    defer (func() {
        if err := recover(); err != nil {
            panic(err)
        }
        checkState(t, result,
                   &expectedState{
                        name: name,
                        result: expectedResult,
                        failed: true,
                        log: log,
                   })
    })()
    result = closure()
    if shouldStop {
        t.Logf("%s didn't stop when it should", name)
    }
}
