/*
Gocheck - A rich testing framework for Go

Copyright (c) 2010, Gustavo Niemeyer <gustavo@niemeyer.net>

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of the copyright holder nor the names of its
      contributors may be used to endorse or promote products derived from
      this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
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
