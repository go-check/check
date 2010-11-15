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

// These tests verify the test running logic.

package gocheck_test

import (
    "gocheck"
    . "gocheck/local"
    "os"
)


var runnerS = gocheck.Suite(&RunS{})

type RunS struct{}

func (s *RunS) TestCountSuite(c *gocheck.C) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Tests ensuring result counting works properly.

func (s *RunS) TestSuccess(c *gocheck.C) {
    output := String{}
    result := gocheck.Run(&SuccessHelper{}, &gocheck.RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 1)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestFailure(c *gocheck.C) {
    output := String{}
    result := gocheck.Run(&FailHelper{}, &gocheck.RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 0)
    c.Check(result.Failed, Equals, 1)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestFixture(c *gocheck.C) {
    output := String{}
    result := gocheck.Run(&FixtureHelper{}, &gocheck.RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 2)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestPanicOnTest(c *gocheck.C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"Test1"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 1)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 1)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestPanicOnSetUpTest(c *gocheck.C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpTest"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 0)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 1)
    c.Check(result.Missed, Equals, 2)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestPanicOnSetUpSuite(c *gocheck.C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpSuite"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 0)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 1)
    c.Check(result.Missed, Equals, 2)
    c.Check(result.RunError, IsNil)
}


// -----------------------------------------------------------------------
// Check result aggregation.

func (s *RunS) TestAdd(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Failed:3, Panicked:4,
                              FixturePanicked:5, Missed:6}
    result.Add(&gocheck.Result{Succeeded:10, Skipped:20, Failed:30,
                               Panicked:40, FixturePanicked:50, Missed:60})
    c.Check(result.Succeeded, Equals, 11)
    c.Check(result.Skipped, Equals, 22)
    c.Check(result.Failed, Equals, 33)
    c.Check(result.Panicked, Equals, 44)
    c.Check(result.FixturePanicked, Equals, 55)
    c.Check(result.Missed, Equals, 66)
    c.Check(result.RunError, IsNil)
}


// -----------------------------------------------------------------------
// Check the Passed() method.

func (s *RunS) TestPassed(c *gocheck.C) {
    c.Assert((&gocheck.Result{}).Passed(), Equals, true)
    c.Assert((&gocheck.Result{Succeeded:1}).Passed(), Equals, true)
    c.Assert((&gocheck.Result{Skipped:1}).Passed(), Equals, true)
    c.Assert((&gocheck.Result{Failed:1}).Passed(), Equals, false)
    c.Assert((&gocheck.Result{Panicked:1}).Passed(), Equals, false)
    c.Assert((&gocheck.Result{FixturePanicked:1}).Passed(), Equals, false)
    c.Assert((&gocheck.Result{Missed:1}).Passed(), Equals, false)
    c.Assert((&gocheck.Result{RunError:os.NewError("!")}).Passed(), Equals, false)
}

// -----------------------------------------------------------------------
// Check that result printing is working correctly.

func (s *RunS) TestPrintSuccess(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:5}
    c.Check(result.String(), Equals, "OK: 5 passed")
}

func (s *RunS) TestPrintFailure(c *gocheck.C) {
    result := &gocheck.Result{Failed:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 FAILED")
}

func (s *RunS) TestPrintSkipped(c *gocheck.C) {
    result := &gocheck.Result{Skipped:5}
    c.Check(result.String(), Equals, "OK: 0 passed, 5 skipped")
}

func (s *RunS) TestPrintPanicked(c *gocheck.C) {
    result := &gocheck.Result{Panicked:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 PANICKED")
}

func (s *RunS) TestPrintFixturePanicked(c *gocheck.C) {
    result := &gocheck.Result{FixturePanicked:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 FIXTURE PANICKED")
}

func (s *RunS) TestPrintMissed(c *gocheck.C) {
    result := &gocheck.Result{Missed:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 MISSED")
}

func (s *RunS) TestPrintAll(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Panicked:3,
                              FixturePanicked:4, Missed:5}
    c.Check(result.String(), Equals,
            "OOPS: 1 passed, 2 skipped, 3 PANICKED, " +
            "4 FIXTURE PANICKED, 5 MISSED")
}

func (s *RunS) TestPrintRunError(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:1, Failed:1,
                              RunError: os.NewError("Kaboom!")}
    c.Check(result.String(), Equals, "ERROR: Kaboom!")
}


// -----------------------------------------------------------------------
// Verify that the method pattern flag works correctly.

func (s *RunS) TestFilterTestName(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "Test[91]"}
    gocheck.Run(&helper, &runConf)
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 5)
}

func (s *RunS) TestFilterTestNameWithAll(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: ".*"}
    gocheck.Run(&helper, &runConf)
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "SetUpTest")
    c.Check(helper.calls[5], Equals, "Test2")
    c.Check(helper.calls[6], Equals, "TearDownTest")
    c.Check(helper.calls[7], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 8)
}

func (s *RunS) TestFilterSuiteName(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "FixtureHelper"}
    gocheck.Run(&helper, &runConf)
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "SetUpTest")
    c.Check(helper.calls[5], Equals, "Test2")
    c.Check(helper.calls[6], Equals, "TearDownTest")
    c.Check(helper.calls[7], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 8)
}

func (s *RunS) TestFilterSuiteNameAndTestName(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "FixtureHelper\\.Test2"}
    gocheck.Run(&helper, &runConf)
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test2")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 5)
}

func (s *RunS) TestFilterAllOut(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "NotFound"}
    gocheck.Run(&helper, &runConf)
    c.Check(helper.n, Equals, 0)
}

func (s *RunS) TestRequirePartialMatch(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "est"}
    gocheck.Run(&helper, &runConf)
    c.Check(helper.n, Equals, 8)
}


func (s *RunS) TestFilterError(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "]["}
    result := gocheck.Run(&helper, &runConf)
    c.Check(result.String(), Equals,
            "ERROR: Bad filter expression: unmatched ']'")
    c.Check(helper.n, Equals, 0)
}
