// Tests for the behavior of the test fixture system.

package gocheck_test


import (
    "gocheck"
    "testing"
)


func TestFixture(t *testing.T) {
    gocheck.RunTestingT(&FixtureS{}, t)
}


// -----------------------------------------------------------------------
// Fixture test suite.

type FixtureS struct{}

func (s *FixtureS) TestOrder(t *gocheck.T) {
    helper := FixtureHelper{}
    gocheck.Run(&helper)
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

func (s *FixtureS) TestPanicOnTest1(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "Test1"}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "SetUpTest")
    t.CheckEqual(helper.calls[5], "Test2")
    t.CheckEqual(helper.calls[6], "TearDownTest")
    t.CheckEqual(helper.calls[7], "TearDownSuite")
    t.CheckEqual(helper.n, 8)

    // t.CheckContains(output, log)
}


// -----------------------------------------------------------------------
// Helper test suite for the fixture tests.

type FixtureHelper struct {
    calls [64]string
    n int
    panicOn string
}

func (s *FixtureHelper) trace(name string) {
    s.calls[s.n] = name
    s.n += 1
    if name == s.panicOn {
        panic(name)
    }
}

func (s *FixtureHelper) SetUpSuite() {
    s.trace("SetUpSuite")
}

func (s *FixtureHelper) TearDownSuite() {
    s.trace("TearDownSuite")
}

func (s *FixtureHelper) SetUpTest() {
    s.trace("SetUpTest")
}

func (s *FixtureHelper) TearDownTest() {
    s.trace("TearDownTest")
}

func (s *FixtureHelper) Test1(t *gocheck.T) {
    s.trace("Test1")
}

func (s *FixtureHelper) Test2(t *gocheck.T) {
    s.trace("Test2")
}


