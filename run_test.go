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
    c.CheckEqual(result.Succeeded, 1)
    c.CheckEqual(result.Failed, 0)
    c.CheckEqual(result.Skipped, 0)
    c.CheckEqual(result.Panicked, 0)
    c.CheckEqual(result.FixturePanicked, 0)
    c.CheckEqual(result.Missed, 0)
    c.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestFailure(c *gocheck.C) {
    output := String{}
    result := gocheck.Run(&FailHelper{}, &gocheck.RunConf{Output: &output})
    c.CheckEqual(result.Succeeded, 0)
    c.CheckEqual(result.Failed, 1)
    c.CheckEqual(result.Skipped, 0)
    c.CheckEqual(result.Panicked, 0)
    c.CheckEqual(result.FixturePanicked, 0)
    c.CheckEqual(result.Missed, 0)
    c.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestFixture(c *gocheck.C) {
    output := String{}
    result := gocheck.Run(&FixtureHelper{}, &gocheck.RunConf{Output: &output})
    c.CheckEqual(result.Succeeded, 2)
    c.CheckEqual(result.Failed, 0)
    c.CheckEqual(result.Skipped, 0)
    c.CheckEqual(result.Panicked, 0)
    c.CheckEqual(result.FixturePanicked, 0)
    c.CheckEqual(result.Missed, 0)
    c.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestPanicOnTest(c *gocheck.C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"Test1"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(result.Succeeded, 1)
    c.CheckEqual(result.Failed, 0)
    c.CheckEqual(result.Skipped, 0)
    c.CheckEqual(result.Panicked, 1)
    c.CheckEqual(result.FixturePanicked, 0)
    c.CheckEqual(result.Missed, 0)
    c.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestPanicOnSetUpTest(c *gocheck.C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpTest"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(result.Succeeded, 0)
    c.CheckEqual(result.Failed, 0)
    c.CheckEqual(result.Skipped, 0)
    c.CheckEqual(result.Panicked, 0)
    c.CheckEqual(result.FixturePanicked, 1)
    c.CheckEqual(result.Missed, 2)
    c.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestPanicOnSetUpSuite(c *gocheck.C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpSuite"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(result.Succeeded, 0)
    c.CheckEqual(result.Failed, 0)
    c.CheckEqual(result.Skipped, 0)
    c.CheckEqual(result.Panicked, 0)
    c.CheckEqual(result.FixturePanicked, 1)
    c.CheckEqual(result.Missed, 2)
    c.CheckEqual(result.RunError, nil)
}


// -----------------------------------------------------------------------
// Check result aggregation.

func (s *RunS) TestAdd(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Failed:3, Panicked:4,
                              FixturePanicked:5, Missed:6}
    result.Add(&gocheck.Result{Succeeded:10, Skipped:20, Failed:30,
                               Panicked:40, FixturePanicked:50, Missed:60})
    c.CheckEqual(result.Succeeded, 11)
    c.CheckEqual(result.Skipped, 22)
    c.CheckEqual(result.Failed, 33)
    c.CheckEqual(result.Panicked, 44)
    c.CheckEqual(result.FixturePanicked, 55)
    c.CheckEqual(result.Missed, 66)
    c.CheckEqual(result.RunError, nil)
}


// -----------------------------------------------------------------------
// Check the Passed() method.

func (s *RunS) TestPassed(c *gocheck.C) {
    c.AssertEqual((&gocheck.Result{}).Passed(), true)
    c.AssertEqual((&gocheck.Result{Succeeded:1}).Passed(), true)
    c.AssertEqual((&gocheck.Result{Skipped:1}).Passed(), true)
    c.AssertEqual((&gocheck.Result{Failed:1}).Passed(), false)
    c.AssertEqual((&gocheck.Result{Panicked:1}).Passed(), false)
    c.AssertEqual((&gocheck.Result{FixturePanicked:1}).Passed(), false)
    c.AssertEqual((&gocheck.Result{Missed:1}).Passed(), false)
    c.AssertEqual((&gocheck.Result{RunError:os.NewError("!")}).Passed(), false)
}

// -----------------------------------------------------------------------
// Check that result printing is working correctly.

func (s *RunS) TestPrintSuccess(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:5}
    c.CheckEqual(result.String(), "OK: 5 passed")
}

func (s *RunS) TestPrintFailure(c *gocheck.C) {
    result := &gocheck.Result{Failed:5}
    c.CheckEqual(result.String(), "OOPS: 0 passed, 5 FAILED")
}

func (s *RunS) TestPrintSkipped(c *gocheck.C) {
    result := &gocheck.Result{Skipped:5}
    c.CheckEqual(result.String(), "OK: 0 passed, 5 skipped")
}

func (s *RunS) TestPrintPanicked(c *gocheck.C) {
    result := &gocheck.Result{Panicked:5}
    c.CheckEqual(result.String(), "OOPS: 0 passed, 5 PANICKED")
}

func (s *RunS) TestPrintFixturePanicked(c *gocheck.C) {
    result := &gocheck.Result{FixturePanicked:5}
    c.CheckEqual(result.String(), "OOPS: 0 passed, 5 FIXTURE PANICKED")
}

func (s *RunS) TestPrintMissed(c *gocheck.C) {
    result := &gocheck.Result{Missed:5}
    c.CheckEqual(result.String(), "OOPS: 0 passed, 5 MISSED")
}

func (s *RunS) TestPrintAll(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Panicked:3,
                              FixturePanicked:4, Missed:5}
    c.CheckEqual(result.String(), "OOPS: 1 passed, 2 skipped, 3 PANICKED, " +
                                  "4 FIXTURE PANICKED, 5 MISSED")
}

func (s *RunS) TestPrintRunError(c *gocheck.C) {
    result := &gocheck.Result{Succeeded:1, Failed:1,
                              RunError: os.NewError("Kaboom!")}
    c.CheckEqual(result.String(), "ERROR: Kaboom!")
}


// -----------------------------------------------------------------------
// Verify that the method pattern flag works correctly.

func (s *RunS) TestFilterTestName(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "Test[91]"}
    gocheck.Run(&helper, &runConf)
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "TearDownSuite")
    c.CheckEqual(helper.n, 5)
}

func (s *RunS) TestFilterTestNameWithAll(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: ".*"}
    gocheck.Run(&helper, &runConf)
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "SetUpTest")
    c.CheckEqual(helper.calls[5], "Test2")
    c.CheckEqual(helper.calls[6], "TearDownTest")
    c.CheckEqual(helper.calls[7], "TearDownSuite")
    c.CheckEqual(helper.n, 8)
}

func (s *RunS) TestFilterSuiteName(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "FixtureHelper"}
    gocheck.Run(&helper, &runConf)
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "SetUpTest")
    c.CheckEqual(helper.calls[5], "Test2")
    c.CheckEqual(helper.calls[6], "TearDownTest")
    c.CheckEqual(helper.calls[7], "TearDownSuite")
    c.CheckEqual(helper.n, 8)
}

func (s *RunS) TestFilterSuiteNameAndTestName(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "FixtureHelper\\.Test2"}
    gocheck.Run(&helper, &runConf)
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test2")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "TearDownSuite")
    c.CheckEqual(helper.n, 5)
}

func (s *RunS) TestFilterAllOut(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "NotFound"}
    gocheck.Run(&helper, &runConf)
    c.CheckEqual(helper.n, 0)
}

func (s *RunS) TestRequirePartialMatch(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "est"}
    gocheck.Run(&helper, &runConf)
    c.CheckEqual(helper.n, 8)
}


func (s *RunS) TestFilterError(c *gocheck.C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "]["}
    result := gocheck.Run(&helper, &runConf)
    c.CheckEqual(result.String(),
                 "ERROR: Bad filter expression: unmatched ']'")
    c.CheckEqual(helper.n, 0)
}
