// These tests verify that the aggregated result counting works correctly.

package gocheck_test

import (
    "gocheck"
)


var resultS = gocheck.Suite(&ResultS{})

type ResultS struct{}

func (s *ResultS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Tests ensuring result counting works properly.

func (s *ResultS) TestSuccess(t *gocheck.T) {
    output := String{}
    result := gocheck.RunWithWriter(&SuccessHelper{}, &output)
    t.CheckEqual(result.Succeeded, 1)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
}

func (s *ResultS) TestFailure(t *gocheck.T) {
    output := String{}
    result := gocheck.RunWithWriter(&FailHelper{}, &output)
    t.CheckEqual(result.Succeeded, 0)
    t.CheckEqual(result.Failed, 1)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
}

func (s *ResultS) TestFixture(t *gocheck.T) {
    output := String{}
    result := gocheck.RunWithWriter(&FixtureHelper{}, &output)
    t.CheckEqual(result.Succeeded, 2)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
}

func (s *ResultS) TestPanicOnTest(t *gocheck.T) {
    output := String{}
    helper := &FixtureHelper{panicOn:"Test1"}
    result := gocheck.RunWithWriter(helper, &output)
    t.CheckEqual(result.Succeeded, 1)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 1)
    t.CheckEqual(result.FixturePanicked, 0)
    t.CheckEqual(result.Missed, 0)
}

func (s *ResultS) TestPanicOnSetUpTest(t *gocheck.T) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpTest"}
    result := gocheck.RunWithWriter(helper, &output)
    t.CheckEqual(result.Succeeded, 0)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 1)
    t.CheckEqual(result.Missed, 2)
}

func (s *ResultS) TestPanicOnSetUpSuite(t *gocheck.T) {
    output := String{}
    helper := &FixtureHelper{panicOn:"SetUpSuite"}
    result := gocheck.RunWithWriter(helper, &output)
    t.CheckEqual(result.Succeeded, 0)
    t.CheckEqual(result.Failed, 0)
    t.CheckEqual(result.Skipped, 0)
    t.CheckEqual(result.Panicked, 0)
    t.CheckEqual(result.FixturePanicked, 1)
    t.CheckEqual(result.Missed, 2)
}


// -----------------------------------------------------------------------
// Check result aggregation.

func (s *ResultS) TestAdd(t *gocheck.T) {
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
}


// -----------------------------------------------------------------------
// Check the Passed() method.

func (s *ResultS) TestPassed(t *gocheck.T) {
    t.AssertEqual((&gocheck.Result{}).Passed(), true)
    t.AssertEqual((&gocheck.Result{Succeeded:1}).Passed(), true)
    t.AssertEqual((&gocheck.Result{Skipped:1}).Passed(), true)
    t.AssertEqual((&gocheck.Result{Failed:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{Panicked:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{FixturePanicked:1}).Passed(), false)
    t.AssertEqual((&gocheck.Result{Missed:1}).Passed(), false)
}

// -----------------------------------------------------------------------
// Check that result printing is working correctly.

func (s *ResultS) TestPrintSuccess(t *gocheck.T) {
    result := &gocheck.Result{Succeeded:5}
    t.CheckEqual(result.String(), "OK: 5 passed")
}

func (s *ResultS) TestPrintFailure(t *gocheck.T) {
    result := &gocheck.Result{Failed:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 FAILED")
}

func (s *ResultS) TestPrintSkipped(t *gocheck.T) {
    result := &gocheck.Result{Skipped:5}
    t.CheckEqual(result.String(), "OK: 0 passed, 5 skipped")
}

func (s *ResultS) TestPrintPanicked(t *gocheck.T) {
    result := &gocheck.Result{Panicked:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 PANICKED")
}

func (s *ResultS) TestPrintFixturePanicked(t *gocheck.T) {
    result := &gocheck.Result{FixturePanicked:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 FIXTURE PANICKED")
}

func (s *ResultS) TestPrintMissed(t *gocheck.T) {
    result := &gocheck.Result{Missed:5}
    t.CheckEqual(result.String(), "OOPS: 0 passed, 5 MISSED")
}

func (s *ResultS) TestPrintAll(t *gocheck.T) {
    result := &gocheck.Result{Succeeded:1, Skipped:2, Panicked:3,
                              FixturePanicked:4, Missed:5}
    t.CheckEqual(result.String(), "OOPS: 1 passed, 2 skipped, 3 PANICKED, " +
                                  "4 FIXTURE PANICKED, 5 MISSED")
}
