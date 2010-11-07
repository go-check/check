package checkers_test

import (
    "testing"
    "gocheck"
    . "gocheck/checkers"
)


func TestGocheck(t *testing.T) {
    gocheck.TestingT(t)
}


type CheckersS struct{}

var _ = gocheck.Suite(&CheckersS{})


func testInfo(c *gocheck.C, checker Checker,
              funcName, obtainedLabel, expectedLabel string) {
    if checker.FuncName() != funcName {
        c.Fatalf("Got function name %s, expected %s",
                 checker.FuncName(), funcName)
    }
    if checker.ObtainedLabel() != obtainedLabel {
        c.Fatalf("Got obtained label %s, expected %s",
                 checker.ObtainedLabel(), obtainedLabel)
    }
    if checker.ExpectedLabel() != expectedLabel {
        c.Fatalf("Got expected label %s, expected %s",
                 checker.ExpectedLabel(), expectedLabel)
    }
}

func testCheck(c *gocheck.C, checkerFunc CheckerFunc,
               obtained, expected interface{}, wantedResult bool) {
    checker := checkerFunc(obtained, expected)
    result := checker.Check()
    if result != wantedResult {
        c.Fatalf("%s(%#v, %#v) returned %#v rather than %#v",
                 checker.FuncName(), obtained, expected, result, wantedResult)
    }
}


func (s *CheckersS) TestEqualsInfo(c *gocheck.C) {
    testInfo(c, Equals(nil, nil), "Equals", "Obtained", "Expected")
}

func (s *CheckersS) TestEqualsCheck(c *gocheck.C) {
    testCheck(c, Equals, 42, 42, true)
    testCheck(c, Equals, 42, 43, false)
    testCheck(c, Equals, int32(42), int64(42), false)
}
