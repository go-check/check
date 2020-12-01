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
	exitCode, output := runHelperSuite("SuccessHelper")
	c.Check(exitCode, Equals, 0)
	c.Check(output.Status("TestLogAndSucceed"), Equals, "PASS")
}

func (s *RunS) TestFailure(c *C) {
	exitCode, output := runHelperSuite("FailHelper")
	c.Check(exitCode, Equals, 1)
	c.Check(output.Status("TestLogAndFail"), Equals, "FAIL")
}

func (s *RunS) TestFixture(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper")
	c.Check(exitCode, Equals, 0)
	c.Check(output.Status("Test1"), Equals, "PASS")
	c.Check(output.Status("Test2"), Equals, "PASS")
}

func (s *RunS) TestPanicOnTest(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper", "-helper.panic", "Test1")
	c.Check(exitCode, Equals, 2)
	c.Check(output.Status("Test1"), Equals, "FAIL")
	// stdlib testing stops on first panic
	c.Check(output.Status("Test2"), Equals, "")
}

func (s *RunS) TestPanicOnSetUpTest(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper", "-helper.panic", "SetUpTest")
	c.Check(exitCode, Equals, 2)
	c.Check(output.Status("Test1"), Equals, "FAIL")
	// stdlib testing stops on first panic
	c.Check(output.Status("Test2"), Equals, "")
}

func (s *RunS) TestPanicOnSetUpSuite(c *C) {
	exitCode, output := runHelperSuite("FixtureHelper", "-helper.panic", "SetUpSuite")
	c.Check(exitCode, Equals, 2)
	// If SetUpSuite fails, no tests from the suite are run
	c.Check(output.Status("Test1"), Equals, "")
	c.Check(output.Status("Test2"), Equals, "")
}

/*
// -----------------------------------------------------------------------
// Verify that List works correctly.

func (s *RunS) TestListFiltered(c *C) {
	names := List(&FixtureHelper{}, &RunConf{Filter: "1"})
	c.Assert(names, DeepEquals, []string{
		"FixtureHelper.Test1",
	})
}

func (s *RunS) TestList(c *C) {
	names := List(&FixtureHelper{}, &RunConf{})
	c.Assert(names, DeepEquals, []string{
		"FixtureHelper.Test1",
		"FixtureHelper.Test2",
	})
}

// -----------------------------------------------------------------------
// Verify that that the keep work dir request indeed does so.

type WorkDirSuite struct {}

func (s *WorkDirSuite) Test(c *C) {
	c.MkDir()
}

func (s *RunS) TestKeepWorkDir(c *C) {
	output := String{}
	runConf := RunConf{Output: &output, Verbose: true, KeepWorkDir: true}
	result := Run(&WorkDirSuite{}, &runConf)

	c.Assert(result.String(), Matches, ".*\nWORK=" + regexp.QuoteMeta(result.WorkDir))

	stat, err := os.Stat(result.WorkDir)
	c.Assert(err, IsNil)
	c.Assert(stat.IsDir(), Equals, true)
}
*/
