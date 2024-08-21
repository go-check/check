// This file contains just a few generic helpers which are used by the
// other test files.

package check_test

import (
	"runtime"
	"testing"

	"gopkg.in/check.v1"
)

// We count the number of suites run at least to get a vague hint that the
// test suite is behaving as it should.  Otherwise a bug introduced at the
// very core of the system could go unperceived.
const suitesRunExpected = 8

var suitesRun int = 0

func Test(t *testing.T) {
	check.TestingT(t)
	if suitesRun != suitesRunExpected {
		t.Fatalf("Expected %d suites to run rather than %d",
			suitesRunExpected, suitesRun)
	}
}

// -----------------------------------------------------------------------
// Helper functions.

// Return the file line where it's called.
func getMyLine() int {
	if _, _, line, ok := runtime.Caller(1); ok {
		return line
	}
	return -1
}

// -----------------------------------------------------------------------
// Helper suite for testing basic fail behavior.

type FailHelper struct{}

func (s *FailHelper) TestLogAndFail(c *check.C) {
	c.Log("Expected failure!")
	c.Fail()
}

// -----------------------------------------------------------------------
// Helper suite for testing basic success behavior.

type SuccessHelper struct{}

func (s *SuccessHelper) TestLogAndSucceed(c *check.C) {
	c.Log("Expected success!")
}

// -----------------------------------------------------------------------
// Helper suite for testing ordering and behavior of fixture.

type FixtureHelper struct {
	calls   int
	failOn  string
	skip    bool
	skipOnN int
}

func (s *FixtureHelper) trace(name string, c *check.C) {
	c.Log(name)
	s.calls += 1
	if name == s.failOn {
		c.FailNow()
	}
	if s.skip && s.skipOnN == s.calls-1 {
		c.Skip("skipOnN == n")
	}
}

func (s *FixtureHelper) SetUpSuite(c *check.C) {
	s.trace("SetUpSuite", c)
}

func (s *FixtureHelper) TearDownSuite(c *check.C) {
	s.trace("TearDownSuite", c)
}

func (s *FixtureHelper) SetUpTest(c *check.C) {
	s.trace("SetUpTest", c)
}

func (s *FixtureHelper) TearDownTest(c *check.C) {
	s.trace("TearDownTest", c)
}

func (s *FixtureHelper) Test1(c *check.C) {
	s.trace("Test1", c)
}

func (s *FixtureHelper) Test2(c *check.C) {
	s.trace("Test2", c)
}
