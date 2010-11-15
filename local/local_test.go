package checkers_test

import (
    "testing"
    "gocheck"
    . "gocheck/local"
    "os"
)


func TestGocheck(t *testing.T) {
    gocheck.TestingT(t)
}


type CheckersS struct{}

var _ = gocheck.Suite(&CheckersS{})


func testInfo(c *gocheck.C, checker Checker,
              name, obtainedVarName, expectedVarName string) {
    if checker.Name() != name {
        c.Fatalf("Got name %s, expected %s", checker.Name(), name)
    }
    obtainedName, expectedName := checker.VarNames()
    if obtainedName != obtainedVarName {
        c.Fatalf("Got obtained label %#v, expected %#v",
                 obtainedName, obtainedVarName)
    }
    if expectedName != expectedVarName {
        c.Fatalf("Got expected label %#v, expected %#v",
                 expectedName, expectedVarName)
    }
}

func testCheck(c *gocheck.C, checker Checker,
               obtained, expected interface{},
               wantedResult bool, wantedError string) {
    result, error := checker.Check(obtained, expected)
    if result != wantedResult || error != wantedError {
        c.Fatalf("%s.Check(%#v, %#v) returned " +
                 "(%#v, %#v) rather than (%#v, %#v)",
                 checker.Name(), obtained, expected,
                 result, error, wantedResult, wantedError)
    }
}

func (s *CheckersS) TestBug(c *gocheck.C) {
    bug := Bug("a %d bc", 42)
    info := bug.GetBugInfo()
    if info != "a 42 bc" {
        c.Fatalf("Bug() returned %#v", info)
    }
}

func (s *CheckersS) TestIsNil(c *gocheck.C) {
    testInfo(c, IsNil, "IsNil", "value", "")

    testCheck(c, IsNil, nil, nil, true, "")
    testCheck(c, IsNil, "a", nil, false, "")

    testCheck(c, IsNil, (chan int)(nil), nil, true, "")
    testCheck(c, IsNil, make(chan int), nil, false, "")
    testCheck(c, IsNil, (os.Error)(nil), nil, true, "")
    testCheck(c, IsNil, os.NewError(""), nil, false, "")
    testCheck(c, IsNil, ([]int)(nil), nil, true, "")
    testCheck(c, IsNil, make([]int, 1), nil, false, "")
}

func (s *CheckersS) TestNotNil(c *gocheck.C) {
    testInfo(c, NotNil, "NotNil", "value", "")

    testCheck(c, NotNil, nil, nil, false, "")
    testCheck(c, NotNil, "a", nil, true, "")

    testCheck(c, NotNil, (chan int)(nil), nil, false, "")
    testCheck(c, NotNil, make(chan int), nil, true, "")
    testCheck(c, NotNil, (os.Error)(nil), nil, false, "")
    testCheck(c, NotNil, os.NewError(""), nil, true, "")
    testCheck(c, NotNil, ([]int)(nil), nil, false, "")
    testCheck(c, NotNil, make([]int, 1), nil, true, "")
}

func (s *CheckersS) TestNot(c *gocheck.C) {
    testInfo(c, Not(IsNil), "Not(IsNil)", "value", "")

    testCheck(c, Not(IsNil), nil, nil, false, "")
    testCheck(c, Not(IsNil), "a", nil, true, "")
}


func (s *CheckersS) TestEquals(c *gocheck.C) {
    testInfo(c, Equals, "Equals", "obtained", "expected")

    // The simplest.
    testCheck(c, Equals, 42, 42, true, "")
    testCheck(c, Equals, 42, 43, false, "")

    // Different native types.
    testCheck(c, Equals, int32(42), int64(42), false, "")

    // With nil.
    testCheck(c, Equals, 42, nil, false, "")

    // Arrays
    testCheck(c, Equals, []byte{1,2}, []byte{1,2}, true, "")
    testCheck(c, Equals, []byte{1,2}, []byte{1,3}, false, "")
}

func (s *CheckersS) TestMatches(c *gocheck.C) {
    testInfo(c, Matches, "Matches", "value", "regex")

    // Simple matching
    testCheck(c, Matches, "abc", "abc", true, "")
    testCheck(c, Matches, "abc", "a.c", true, "")

    // Must match fully
    testCheck(c, Matches, "abc", "ab", false, "")
    testCheck(c, Matches, "abc", "bc", false, "")

    // String()-enabled values accepted
    testCheck(c, Matches, os.NewError("abc"), "a.c", true, "")
    testCheck(c, Matches, os.NewError("abc"), "a.d", false, "")

    // Some error conditions.
    testCheck(c, Matches, 1, "a.c", false,
              "Obtained value is not a string and has no .String()")
    testCheck(c, Matches, "abc", "a[c", false,
              "Can't compile regex: unmatched '['")
}
