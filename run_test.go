// These tests verify the test running logic.

package gocheck_test

import (
    . "gocheck"
    "os"
)


var runnerS = Suite(&RunS{})

type RunS struct{}

func (s *RunS) TestCountSuite(c *C) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Tests ensuring result counting works properly.

func (s *RunS) TestSuccess(c *C) {
    output := String{}
    result := Run(&SuccessHelper{}, &RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 1)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestFailure(c *C) {
    output := String{}
    result := Run(&FailHelper{}, &RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 0)
    c.Check(result.Failed, Equals, 1)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestFixture(c *C) {
    output := String{}
    result := Run(&FixtureHelper{}, &RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 2)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestPanicOnTest(c *C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"Test1"}
    result := Run(helper, &RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 1)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 1)
    c.Check(result.FixturePanicked, Equals, 0)
    c.Check(result.Missed, Equals, 0)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestPanicOnSetUpTest(c *C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpTest"}
    result := Run(helper, &RunConf{Output: &output})
    c.Check(result.Succeeded, Equals, 0)
    c.Check(result.Failed, Equals, 0)
    c.Check(result.Skipped, Equals, 0)
    c.Check(result.Panicked, Equals, 0)
    c.Check(result.FixturePanicked, Equals, 1)
    c.Check(result.Missed, Equals, 2)
    c.Check(result.RunError, IsNil)
}

func (s *RunS) TestPanicOnSetUpSuite(c *C) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpSuite"}
    result := Run(helper, &RunConf{Output: &output})
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

func (s *RunS) TestAdd(c *C) {
    result := &Result{Succeeded:1, Skipped:2, Failed:3, Panicked:4,
                              FixturePanicked:5, Missed:6}
    result.Add(&Result{Succeeded:10, Skipped:20, Failed:30,
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

func (s *RunS) TestPassed(c *C) {
    c.Assert((&Result{}).Passed(), Equals, true)
    c.Assert((&Result{Succeeded:1}).Passed(), Equals, true)
    c.Assert((&Result{Skipped:1}).Passed(), Equals, true)
    c.Assert((&Result{Failed:1}).Passed(), Equals, false)
    c.Assert((&Result{Panicked:1}).Passed(), Equals, false)
    c.Assert((&Result{FixturePanicked:1}).Passed(), Equals, false)
    c.Assert((&Result{Missed:1}).Passed(), Equals, false)
    c.Assert((&Result{RunError:os.NewError("!")}).Passed(), Equals, false)
}

// -----------------------------------------------------------------------
// Check that result printing is working correctly.

func (s *RunS) TestPrintSuccess(c *C) {
    result := &Result{Succeeded:5}
    c.Check(result.String(), Equals, "OK: 5 passed")
}

func (s *RunS) TestPrintFailure(c *C) {
    result := &Result{Failed:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 FAILED")
}

func (s *RunS) TestPrintSkipped(c *C) {
    result := &Result{Skipped:5}
    c.Check(result.String(), Equals, "OK: 0 passed, 5 skipped")
}

func (s *RunS) TestPrintPanicked(c *C) {
    result := &Result{Panicked:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 PANICKED")
}

func (s *RunS) TestPrintFixturePanicked(c *C) {
    result := &Result{FixturePanicked:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 FIXTURE PANICKED")
}

func (s *RunS) TestPrintMissed(c *C) {
    result := &Result{Missed:5}
    c.Check(result.String(), Equals, "OOPS: 0 passed, 5 MISSED")
}

func (s *RunS) TestPrintAll(c *C) {
    result := &Result{Succeeded:1, Skipped:2, Panicked:3,
                              FixturePanicked:4, Missed:5}
    c.Check(result.String(), Equals,
            "OOPS: 1 passed, 2 skipped, 3 PANICKED, " +
            "4 FIXTURE PANICKED, 5 MISSED")
}

func (s *RunS) TestPrintRunError(c *C) {
    result := &Result{Succeeded:1, Failed:1,
                              RunError: os.NewError("Kaboom!")}
    c.Check(result.String(), Equals, "ERROR: Kaboom!")
}


// -----------------------------------------------------------------------
// Verify that the method pattern flag works correctly.

func (s *RunS) TestFilterTestName(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: "Test[91]"}
    Run(&helper, &runConf)
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 5)
}

func (s *RunS) TestFilterTestNameWithAll(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: ".*"}
    Run(&helper, &runConf)
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

func (s *RunS) TestFilterSuiteName(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: "FixtureHelper"}
    Run(&helper, &runConf)
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

func (s *RunS) TestFilterSuiteNameAndTestName(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: "FixtureHelper\\.Test2"}
    Run(&helper, &runConf)
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test2")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 5)
}

func (s *RunS) TestFilterAllOut(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: "NotFound"}
    Run(&helper, &runConf)
    c.Check(helper.n, Equals, 0)
}

func (s *RunS) TestRequirePartialMatch(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: "est"}
    Run(&helper, &runConf)
    c.Check(helper.n, Equals, 8)
}


func (s *RunS) TestFilterError(c *C) {
    helper := FixtureHelper{}
    output := String{}
    runConf := RunConf{Output: &output, Filter: "]["}
    result := Run(&helper, &runConf)
    c.Check(result.String(), Equals,
            "ERROR: Bad filter expression: unmatched ']'")
    c.Check(helper.n, Equals, 0)
}
