// These tests verify the test running logic.

package check_test

import (
	. "gopkg.in/check.v1"
)

var runnerS = Suite(&RunS{})

type RunS struct{}

func (s *RunS) TestCountSuite(c *C) {
	suitesRun += 1
}

// -----------------------------------------------------------------------
// Tests ensuring result counting works properly.

func (s *RunS) TestSuccess(c *C) {
	exitCode, output := runHelperSuite(c, "SuccessHelper")
	c.Check(exitCode, Equals, 0)
	c.Check(output.Status("SuccessHelper/TestLogAndSucceed"), Equals, "pass")
}

func (s *RunS) TestFailure(c *C) {
	exitCode, output := runHelperSuite(c, "FailHelper")
	c.Check(exitCode, Equals, 1)
	c.Check(output.Status("FailHelper/TestLogAndFail"), Equals, "fail")
}

func (s *RunS) TestFixture(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper")
	c.Check(exitCode, Equals, 0)
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "pass")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "pass")
}

func (s *RunS) TestFailOnTest(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "Test1")
	c.Check(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "pass")
}

func (s *RunS) TestFailOnSetUpTest(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "SetUpTest")
	c.Check(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "fail")
}

func (s *RunS) TestFailOnSetUpSuite(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "SetUpSuite")
	c.Check(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	// If SetUpSuite fails, no tests from the suite are run
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "")
}
