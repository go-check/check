// Tests for the behavior of the test fixture system.

package check_test

import (
	. "gopkg.in/check.v1"
)

// -----------------------------------------------------------------------
// Fixture test suite.

type FixtureS struct{}

var fixtureS = Suite(&FixtureS{})

func (s *FixtureS) TestCountSuite(c *C) {
	suitesRun += 1
}

// -----------------------------------------------------------------------
// Basic fixture ordering verification.

func (s *FixtureS) TestOrder(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper")
	c.Assert(exitCode, Equals, 0)
	c.Check(output.Status("FixtureHelper"), Equals, "pass")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test1\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test2\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

// -----------------------------------------------------------------------
// Check the behavior when fatal errors occur within tests and fixtures.

func (s *FixtureS) TestFailOnTest(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "Test1")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "pass")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test1\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test2\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

func (s *FixtureS) TestFailOnSetUpTest(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "SetUpTest")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "fail")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

func (s *FixtureS) TestFailOnTearDownTest(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "TearDownTest")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "fail")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test1\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test2\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

func (s *FixtureS) TestFailOnSetUpSuite(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "SetUpSuite")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

func (s *FixtureS) TestFailOnTearDownSuite(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.fail", "TearDownSuite")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureHelper"), Equals, "fail")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "pass")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "pass")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test1\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test2\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

// -----------------------------------------------------------------------
// A wrong argument on a test or fixture will produce a nice error.

func (s *FixtureS) TestFailOnWrongTestArg(c *C) {
	exitCode, output := runHelperSuite(c, "WrongTestArgHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("WrongTestArgHelper/Test1"), Equals, "fail")
	c.Check(output.Logs("WrongTestArgHelper/Test1"), Matches, ""+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check.go:\\d+: bad signature for method Test1: func\\(int\\)\n"+
		"    check_test.go:\\d+: TearDownTest\n")
}

func (s *FixtureS) TestFailOnWrongSetUpTestArg(c *C) {
	exitCode, output := runHelperSuite(c, "WrongSetUpTestArgHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("WrongSetUpTestArgHelper/Test1"), Equals, "fail")
	c.Check(output.Logs("WrongSetUpTestArgHelper/Test1"), Matches, ""+
		"    check.go:\\d+: bad signature for method SetUpTest: func\\(int\\)\n"+
		"    check_test.go:\\d+: TearDownTest\n")
}

func (s *FixtureS) TestFailOnWrongSetUpSuiteArg(c *C) {
	exitCode, output := runHelperSuite(c, "WrongSetUpSuiteArgHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("WrongSetUpSuiteArgHelper"), Equals, "fail")
	c.Check(output.Logs("WrongSetUpSuiteArgHelper"), Matches, ""+
		"    check.go:\\d+: bad signature for method SetUpSuite: func\\(int\\)\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
	c.Check(output.Status("WrongSetUpSuiteArgHelper/Test1"), Equals, "")
}

// -----------------------------------------------------------------------
// Nice errors also when tests or fixture have wrong arg count.

func (s *FixtureS) TestPanicOnWrongTestArgCount(c *C) {
	exitCode, output := runHelperSuite(c, "WrongTestArgCountHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("WrongTestArgCountHelper/Test1"), Equals, "fail")
	c.Check(output.Logs("WrongTestArgCountHelper/Test1"), Matches, ""+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check.go:\\d+: bad signature for method Test1: func\\(\\*check.C, int\\)\n"+
		"    check_test.go:\\d+: TearDownTest\n")
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArgCount(c *C) {
	exitCode, output := runHelperSuite(c, "WrongSetUpTestArgCountHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("WrongSetUpTestArgCountHelper/Test1"), Equals, "fail")
	c.Check(output.Logs("WrongSetUpTestArgCountHelper/Test1"), Matches, ""+
		"    check.go:\\d+: bad signature for method SetUpTest: func\\(\\*check.C, int\\)\n"+
		"    check_test.go:\\d+: TearDownTest\n")
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArgCount(c *C) {
	exitCode, output := runHelperSuite(c, "WrongSetUpSuiteArgCountHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("WrongSetUpSuiteArgCountHelper"), Equals, "fail")
	c.Check(output.Logs("WrongSetUpSuiteArgCountHelper"), Matches, ""+
		"    check.go:\\d+: bad signature for method SetUpSuite: func\\(\\*check.C, int\\)\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
	c.Check(output.Status("WrongSetUpSuiteArgCountHelper/Test1"), Equals, "")
}

// -----------------------------------------------------------------------
// Helper test suites with wrong function arguments.

type WrongTestArgHelper struct {
	FixtureHelper
}

func (s *WrongTestArgHelper) Test1(t int) {
}

type WrongSetUpTestArgHelper struct {
	FixtureHelper
}

func (s *WrongSetUpTestArgHelper) SetUpTest(t int) {
}

type WrongSetUpSuiteArgHelper struct {
	FixtureHelper
}

func (s *WrongSetUpSuiteArgHelper) SetUpSuite(t int) {
}

type WrongTestArgCountHelper struct {
	FixtureHelper
}

func (s *WrongTestArgCountHelper) Test1(c *C, i int) {
}

type WrongSetUpTestArgCountHelper struct {
	FixtureHelper
}

func (s *WrongSetUpTestArgCountHelper) SetUpTest(c *C, i int) {
}

type WrongSetUpSuiteArgCountHelper struct {
	FixtureHelper
}

func (s *WrongSetUpSuiteArgCountHelper) SetUpSuite(c *C, i int) {
}

// -----------------------------------------------------------------------
// Ensure fixture doesn't run without tests.

type NoTestsHelper struct{}

func (s *NoTestsHelper) SetUpSuite(c *C) {
	c.Fatal("SetUpSuite called")
}

func (s *NoTestsHelper) TearDownSuite(c *C) {
	c.Fatal("TearDownSuite called")
}

func (s *FixtureS) TestFixtureDoesntRunWithoutTests(c *C) {
	exitCode, output := runHelperSuite(c, "NoTestsHelper")
	c.Assert(exitCode, Equals, 0)
	c.Check(output.Status("NoTestsHelper"), Equals, "pass")
}

// -----------------------------------------------------------------------
// Verify that checks and assertions work correctly inside the fixture.

type FixtureCheckHelper struct {
	fail      string
}

func (s *FixtureCheckHelper) SetUpSuite(c *C) {
	switch s.fail {
	case "SetUpSuiteAssert":
		c.Assert(false, Equals, true)
	case "SetUpSuiteCheck":
		c.Check(false, Equals, true)
	}
}

func (s *FixtureCheckHelper) SetUpTest(c *C) {
	switch s.fail {
	case "SetUpTestAssert":
		c.Assert(false, Equals, true)
	case "SetUpTestCheck":
		c.Check(false, Equals, true)
	}
}

func (s *FixtureCheckHelper) Test(c *C) {
	// Do nothing.
}

func (s *FixtureS) TestSetUpSuiteCheck(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureCheckHelper", "-helper.fail", "SetUpSuiteCheck")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureCheckHelper"), Equals, "fail")
	c.Check(output.Status("FixtureCheckHelper/Test"), Equals, "pass")
	c.Check(output.Logs("FixtureCheckHelper"), Matches, ""+
		"    fixture_test.go:\\d+: \n"+
		"            c.Check\\(false, Equals, true\\)\n"+
		"        ... obtained bool = false\n"+
		"        ... expected bool = true\n")
}

func (s *FixtureS) TestSetUpSuiteAssert(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureCheckHelper", "-helper.fail", "SetUpSuiteAssert")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureCheckHelper"), Equals, "fail")
	c.Check(output.Status("FixtureCheckHelper/Test"), Equals, "")
	c.Check(output.Logs("FixtureCheckHelper"), Matches, ""+
		"    fixture_test.go:\\d+: \n"+
		"            c.Assert\\(false, Equals, true\\)\n"+
		"        ... obtained bool = false\n"+
		"        ... expected bool = true\n")
}

// -----------------------------------------------------------------------
// Verify that logging within SetUpTest() persists within the test log itself.

type FixtureLogHelper struct {
	c *C
}

func (s *FixtureLogHelper) SetUpTest(c *C) {
	s.c = c
	c.Log("1")
}

func (s *FixtureLogHelper) Test(c *C) {
	c.Log("2")
	s.c.Log("3")
	c.Log("4")
	c.Fail()
}

func (s *FixtureLogHelper) TearDownTest(c *C) {
	s.c.Log("5")
}

func (s *FixtureS) TestFixtureLogging(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureLogHelper")
	c.Assert(exitCode, Equals, 1)
	c.Check(output.Status("FixtureLogHelper/Test"), Equals, "fail")
	c.Check(output.Logs("FixtureLogHelper/Test"), Matches, ""+
		"    fixture_test.go:\\d+: 1\n"+
		"    fixture_test.go:\\d+: 2\n"+
		"    fixture_test.go:\\d+: 3\n"+
		"    fixture_test.go:\\d+: 4\n"+
		"    fixture_test.go:\\d+: 5\n")
}

// -----------------------------------------------------------------------
// Skip() within fixture methods.

func (s *FixtureS) TestSkipSuite(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.skip", "0")
	c.Assert(exitCode, Equals, 0)
	c.Check(output.Status("FixtureHelper"), Equals, "skip")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: skipOnN == n\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}

func (s *FixtureS) TestSkipTest(c *C) {
	exitCode, output := runHelperSuite(c, "FixtureHelper", "-helper.skip", "1")
	c.Assert(exitCode, Equals, 0)
	c.Check(output.Status("FixtureHelper"), Equals, "pass")
	c.Check(output.Status("FixtureHelper/Test1"), Equals, "skip")
	c.Check(output.Status("FixtureHelper/Test2"), Equals, "pass")
	c.Check(output.AllLogs(), Matches, ""+
		"    check_test.go:\\d+: SetUpSuite\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: skipOnN == n\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: SetUpTest\n"+
		"    check_test.go:\\d+: Test2\n"+
		"    check_test.go:\\d+: TearDownTest\n"+
		"    check_test.go:\\d+: TearDownSuite\n")
}
