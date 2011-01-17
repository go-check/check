package gocheck_test

import (
    "gocheck"
    "os"
)


type CheckersS struct{}

var _ = gocheck.Suite(&CheckersS{})


func testInfo(c *gocheck.C, checker gocheck.Checker,
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

func testCheck(c *gocheck.C, checker gocheck.Checker,
obtained, expected interface{},
wantedResult bool, wantedError string) {
    result, error := checker.Check(obtained, expected)
    if result != wantedResult || error != wantedError {
        c.Fatalf("%s.Check(%#v, %#v) returned "+
            "(%#v, %#v) rather than (%#v, %#v)",
            checker.Name(), obtained, expected,
            result, error, wantedResult, wantedError)
    }
}

func (s *CheckersS) TestBug(c *gocheck.C) {
    bug := gocheck.Bug("a %d bc", 42)
    info := bug.GetBugInfo()
    if info != "a 42 bc" {
        c.Fatalf("Bug() returned %#v", info)
    }
}

func (s *CheckersS) TestIsNil(c *gocheck.C) {
    testInfo(c, gocheck.IsNil, "IsNil", "value", "")

    testCheck(c, gocheck.IsNil, nil, nil, true, "")
    testCheck(c, gocheck.IsNil, "a", nil, false, "")

    testCheck(c, gocheck.IsNil, (chan int)(nil), nil, true, "")
    testCheck(c, gocheck.IsNil, make(chan int), nil, false, "")
    testCheck(c, gocheck.IsNil, (os.Error)(nil), nil, true, "")
    testCheck(c, gocheck.IsNil, os.NewError(""), nil, false, "")
    testCheck(c, gocheck.IsNil, ([]int)(nil), nil, true, "")
    testCheck(c, gocheck.IsNil, make([]int, 1), nil, false, "")
    testCheck(c, gocheck.IsNil, int(0), nil, false, "")
}

func (s *CheckersS) TestNotNil(c *gocheck.C) {
    testInfo(c, gocheck.NotNil, "NotNil", "value", "")

    testCheck(c, gocheck.NotNil, nil, nil, false, "")
    testCheck(c, gocheck.NotNil, "a", nil, true, "")

    testCheck(c, gocheck.NotNil, (chan int)(nil), nil, false, "")
    testCheck(c, gocheck.NotNil, make(chan int), nil, true, "")
    testCheck(c, gocheck.NotNil, (os.Error)(nil), nil, false, "")
    testCheck(c, gocheck.NotNil, os.NewError(""), nil, true, "")
    testCheck(c, gocheck.NotNil, ([]int)(nil), nil, false, "")
    testCheck(c, gocheck.NotNil, make([]int, 1), nil, true, "")
}

func (s *CheckersS) TestNot(c *gocheck.C) {
    testInfo(c, gocheck.Not(gocheck.IsNil), "Not(IsNil)", "value", "")

    testCheck(c, gocheck.Not(gocheck.IsNil), nil, nil, false, "")
    testCheck(c, gocheck.Not(gocheck.IsNil), "a", nil, true, "")
}


type simpleStruct struct {
    i int
}

func (s *CheckersS) TestEquals(c *gocheck.C) {
    testInfo(c, gocheck.Equals, "Equals", "obtained", "expected")

    // The simplest.
    testCheck(c, gocheck.Equals, 42, 42, true, "")
    testCheck(c, gocheck.Equals, 42, 43, false, "")

    // Different native types.
    testCheck(c, gocheck.Equals, int32(42), int64(42), false, "")

    // With nil.
    testCheck(c, gocheck.Equals, 42, nil, false, "")

    // Arrays
    testCheck(c, gocheck.Equals, []byte{1, 2}, []byte{1, 2}, true, "")
    testCheck(c, gocheck.Equals, []byte{1, 2}, []byte{1, 3}, false, "")

    // Struct values
    testCheck(c, gocheck.Equals, simpleStruct{1}, simpleStruct{1}, true, "")
    testCheck(c, gocheck.Equals, simpleStruct{1}, simpleStruct{2}, false, "")

    // Struct pointers
    testCheck(c, gocheck.Equals, &simpleStruct{1}, &simpleStruct{1}, true, "")
    testCheck(c, gocheck.Equals, &simpleStruct{1}, &simpleStruct{2}, false, "")
}

func (s *CheckersS) TestMatches(c *gocheck.C) {
    testInfo(c, gocheck.Matches, "Matches", "value", "regex")

    // Simple matching
    testCheck(c, gocheck.Matches, "abc", "abc", true, "")
    testCheck(c, gocheck.Matches, "abc", "a.c", true, "")

    // Must match fully
    testCheck(c, gocheck.Matches, "abc", "ab", false, "")
    testCheck(c, gocheck.Matches, "abc", "bc", false, "")

    // String()-enabled values accepted
    testCheck(c, gocheck.Matches, os.NewError("abc"), "a.c", true, "")
    testCheck(c, gocheck.Matches, os.NewError("abc"), "a.d", false, "")

    // Some error conditions.
    testCheck(c, gocheck.Matches, 1, "a.c", false,
        "Obtained value is not a string and has no .String()")
    testCheck(c, gocheck.Matches, "abc", "a[c", false,
        "Can't compile regex: unmatched '['")
}
