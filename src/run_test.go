// These tests verify the test running logic.

package gocheck_test

import (
    "gocheck"
    "os"
)


var runnerS = gocheck.Suite(&RunS{})

type RunS struct{}

func (s *RunS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Tests ensuring result counting works properly.

func (s *RunS) TestSuccess(t *gocheck.T) {
    output := String{}
    result := gocheck.Run(&SuccessHelper{}, &gocheck.RunConf{Output: &output})
    t.CheckEqual(result.Succeeded, 1)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
    t.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestFailure(t *gocheck.T) {
    output := String{}
    result := gocheck.Run(&FailHelper{}, &gocheck.RunConf{Output: &output})
    t.CheckEqual(result.Succeeded, 0)
    t.CheckEqual(result.Failed, 1)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
    t.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestFixture(t *gocheck.T) {
    output := String{}
    result := gocheck.Run(&FixtureHelper{}, &gocheck.RunConf{Output: &output})
    t.CheckEqual(result.Succeeded, 2)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
    t.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestPanicOnTest(t *gocheck.T) {
    output := String{}
    helper := &FixtureHelper{panicOn:"Test1"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(result.Succeeded, 1)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 1)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
    t.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestPanicOnSetUpTest(t *gocheck.T) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpTest"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(result.Succeeded, 0)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 1)
    t.CheckEqual(result.Missed, 2)
    t.CheckEqual(result.RunError, nil)
}

func (s *RunS) TestPanicOnSetUpSuite(t *gocheck.T) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpSuite"}
    result := gocheck.Run(helper, &gocheck.RunConf{Output: &output})
    t.CheckEqual(result.Succeeded, 0)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 1)
    t.CheckEqual(result.Missed, 2)
    t.CheckEqual(result.RunError, nil)
}


// -----------------------------------------------------------------------
// Check result aggregation.

func (s *RunS) TestAdd(t *gocheck.T) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Failed:3, Panicked:4,
                              FixturePanicked:5, Missed:6}
    result.Add(&gocheck.Result{Succeeded:10, Skipped:20, Failed:30,
                               Panicked:40, FixturePanicked:50, Missed:60})
    t.CheckEqual(result.Succeeded, 11)
    t.CheckEqual(result.Skipped, 22)
    t.CheckEqual(result.Failed, 33)
    t.CheckEqual(result.Panicked, 44)
    t.CheckEqual(result.FixturePanicked, 55)
    t.CheckEqual(result.Missed, 66)
    t.CheckEqual(result.RunError, nil)
}


// -----------------------------------------------------------------------
// Check the Passed() method.

func (s *RunS) TestPassed(t *gocheck.T) {
    t.AssertEqual((&gocheck.Result{}).Passed(), true)
    t.AssertEqual((&gocheck.Result{Succeeded:1}).Passed(), true)
    t.AssertEqual((&gocheck.Result{Skipped:1}).Passed(), true)
    t.AssertEqual((&gocheck.Result{Failed:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{Panicked:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{FixturePanicked:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{Missed:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{RunError:os.NewError("!")}).Passed(), false)
}

// -----------------------------------------------------------------------
// Check that result printing is working correctly.

func (s *RunS) TestPrintSuccess(t *gocheck.T) {
    result := &gocheck.Result{Succeeded:5}
    t.CheckEqual(result.String(), "OK: 5 passed")
}

func (s *RunS) TestPrintFailure(t *gocheck.T) {
    result := &gocheck.Result{Failed:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 FAILED")
}

func (s *RunS) TestPrintSkipped(t *gocheck.T) {
    result := &gocheck.Result{Skipped:5}
    t.CheckEqual(result.String(), "OK: 0 passed, 5 skipped")
}

func (s *RunS) TestPrintPanicked(t *gocheck.T) {
    result := &gocheck.Result{Panicked:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 PANICKED")
}

func (s *RunS) TestPrintFixturePanicked(t *gocheck.T) {
    result := &gocheck.Result{FixturePanicked:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 FIXTURE PANICKED")
}

func (s *RunS) TestPrintMissed(t *gocheck.T) {
    result := &gocheck.Result{Missed:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 MISSED")
}

func (s *RunS) TestPrintAll(t *gocheck.T) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Panicked:3,
                              FixturePanicked:4, Missed:5}
    t.CheckEqual(result.String(), "OOPS: 1 passed, 2 skipped, 3 PANICKED, " +
                                  "4 FIXTURE PANICKED, 5 MISSED")
}

func (s *RunS) TestPrintRunError(t *gocheck.T) {
    result := &gocheck.Result{Succeeded:1, Failed:1,
                              RunError: os.NewError("Kaboom!")}
    t.CheckEqual(result.String(), "ERROR: Kaboom!")
}


// -----------------------------------------------------------------------
// Verify that the method pattern flag works correctly.

func (s *RunS) TestFilterTestName(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "Test[91]"}
    gocheck.Run(&helper, &runConf)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "TearDownSuite")
    t.CheckEqual(helper.n, 5)
}

func (s *RunS) TestFilterTestNameWithAll(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: ".*"}
    gocheck.Run(&helper, &runConf)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "SetUpTest")
    t.CheckEqual(helper.calls[5], "Test2")
    t.CheckEqual(helper.calls[6], "TearDownTest")
    t.CheckEqual(helper.calls[7], "TearDownSuite")
    t.CheckEqual(helper.n, 8)
}

func (s *RunS) TestFilterSuiteName(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "FixtureHelper"}
    gocheck.Run(&helper, &runConf)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "SetUpTest")
    t.CheckEqual(helper.calls[5], "Test2")
    t.CheckEqual(helper.calls[6], "TearDownTest")
    t.CheckEqual(helper.calls[7], "TearDownSuite")
    t.CheckEqual(helper.n, 8)
}

func (s *RunS) TestFilterSuiteNameAndTestName(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "FixtureHelper\\.Test2"}
    gocheck.Run(&helper, &runConf)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test2")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "TearDownSuite")
    t.CheckEqual(helper.n, 5)
}

func (s *RunS) TestFilterAllOut(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "NotFound"}
    gocheck.Run(&helper, &runConf)
    t.CheckEqual(helper.n, 0)
}

func (s *RunS) TestRequireFullMatch(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "Test"}
    gocheck.Run(&helper, &runConf)
    t.CheckEqual(helper.n, 0)
}


func (s *RunS) TestFilterError(t *gocheck.T) {
    helper := FixtureHelper{}
    output := String{}
    runConf := gocheck.RunConf{Output: &output, Filter: "]["}
    result := gocheck.Run(&helper, &runConf)
    t.CheckEqual(result.String(),
                 "ERROR: Bad filter expression: unmatched ']'")
    t.CheckEqual(helper.n, 0)
}
