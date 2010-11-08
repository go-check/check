package checkers_test

import (
    "testing"
    "gocheck"
    . "gocheck/checkers"
    "os"
)


func TestGocheck(t *testing.T) {
    gocheck.TestingT(t)
}


type CheckersS struct{}

var _ = gocheck.Suite(&CheckersS{})


func testInfo(c *gocheck.C, checker Checker, name,
              obtainedVarName, obtainedVarLabel,
              expectedVarName, expectedVarLabel string) {
    if checker.Name() != name {
        c.Fatalf("Got name %s, expected %s", checker.Name(), name)
    }
    varName, varLabel := checker.ObtainedLabel()
    if varName != obtainedVarName || varLabel != obtainedVarLabel {
        c.Fatalf("Got obtained label (%#v, %#v), expected (%#v, %#v)",
                 varName, varLabel, obtainedVarName, obtainedVarLabel)
    }
    varName, varLabel = checker.ExpectedLabel()
    if varName != expectedVarName || varLabel != expectedVarLabel {
        c.Fatalf("Got expected label (%#v, %#v), expected (%#v, %#v)",
                 varName, varLabel, expectedVarName, expectedVarLabel)
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


func (s *CheckersS) TestEquals(c *gocheck.C) {
    testInfo(c, Equals, "Equals",
             "obtained", "Obtained", "expected", "Expected")

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
    testInfo(c, Matches, "Matches",
             "value", "Value", "regex", "Expected to match regex")

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
