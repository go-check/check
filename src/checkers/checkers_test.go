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
              name, obtainedLabel, expectedLabel string) {
    if checker.Name() != name {
        c.Fatalf("Got name %s, expected %s", checker.Name(), name)
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

func testCheck(c *gocheck.C, checker Checker,
               obtained, expected interface{}, wantedResult bool) {
    result := checker.Check(obtained, expected)
    if result != wantedResult {
        c.Fatalf("%s.Check(%#v, %#v) returned %#v rather than %#v",
                 checker.Name(), obtained, expected, result, wantedResult)
    }
}


func (s *CheckersS) TestEqualsInfo(c *gocheck.C) {
    testInfo(c, Equals, "Equals", "Obtained", "Expected")
}

func (s *CheckersS) TestEqualsCheck(c *gocheck.C) {
    testCheck(c, Equals, 42, 42, true)
    testCheck(c, Equals, 42, 43, false)
    testCheck(c, Equals, int32(42), int64(42), false)
}
