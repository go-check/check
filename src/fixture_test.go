// Tests for the behavior of the test fixture system.

package gocheck_test


import (
    "gocheck"
    "testing"
)


// -----------------------------------------------------------------------
// Fixture test suite.

type FixtureS struct{}

var fixtureS = gocheck.Suite(&FixtureS{})

func (s *FixtureS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}

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

    expected := "^\n-+\n" +
                "FAIL: fixture_test\\.go:FixtureHelper.Test1\n\n" +
                "\\.\\.\\. Panic: Test1 \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.Test1\n$"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
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


