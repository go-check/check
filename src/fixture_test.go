// Tests for the behavior of the test fixture system.

package gocheck_test


import (
    "gocheck"
    . "gocheck/local"
)


// -----------------------------------------------------------------------
// Fixture test suite.

type FixtureS struct{}

var fixtureS = gocheck.Suite(&FixtureS{})

func (s *FixtureS) TestCountSuite(c *gocheck.C) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Basic fixture ordering verification.

func (s *FixtureS) TestOrder(c *gocheck.C) {
    helper := FixtureHelper{}
    gocheck.Run(&helper, nil)
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


// -----------------------------------------------------------------------
// Check the behavior when panics occur within tests and fixtures.

func (s *FixtureS) TestPanicOnTest(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "Test1"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "SetUpTest")
    c.Check(helper.calls[5], Equals, "Test2")
    c.Check(helper.calls[6], Equals, "TearDownTest")
    c.Check(helper.calls[7], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 8)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: FixtureHelper.Test1\n\n" +
                "\\.\\.\\. Panic: Test1 \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.Test1\n$"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnSetUpTest(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "SetUpTest"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "TearDownTest")
    c.Check(helper.calls[3], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 4)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper\\.SetUpTest\n\n" +
                "\\.\\.\\. Panic: SetUpTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpTest\n" +
                "\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: Fixture has panicked " +
                "\\(see related PANIC\\)\n$"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnTearDownTest(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "TearDownTest"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 5)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper.TearDownTest\n\n" +
                "\\.\\.\\. Panic: TearDownTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownTest\n" +
                "\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: Fixture has panicked " +
                "\\(see related PANIC\\)\n$"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnSetUpSuite(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "SetUpSuite"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 2)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper.SetUpSuite\n\n" +
                "\\.\\.\\. Panic: SetUpSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpSuite\n$"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnTearDownSuite(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "TearDownSuite"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "Test1")
    c.Check(helper.calls[3], Equals, "TearDownTest")
    c.Check(helper.calls[4], Equals, "SetUpTest")
    c.Check(helper.calls[5], Equals, "Test2")
    c.Check(helper.calls[6], Equals, "TearDownTest")
    c.Check(helper.calls[7], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 8)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: " +
                "FixtureHelper.TearDownSuite\n\n" +
                "\\.\\.\\. Panic: TearDownSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownSuite\n$"

    c.Check(output.value, Matches, expected)
}


// -----------------------------------------------------------------------
// A wrong argument on a test or fixture will produce a nice error.

func (s *FixtureS) TestPanicOnWrongTestArg(c *gocheck.C) {
    helper := WrongTestArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "TearDownTest")
    c.Check(helper.calls[3], Equals, "SetUpTest")
    c.Check(helper.calls[4], Equals, "Test2")
    c.Check(helper.calls[5], Equals, "TearDownTest")
    c.Check(helper.calls[6], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:[0-9]+: " +
                "WrongTestArgHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgHelper\\.Test1 argument " +
                "should be \\*gocheck\\.C\n"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArg(c *gocheck.C) {
    helper := WrongSetUpTestArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.n, Equals, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpTestArgHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.C\n"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArg(c *gocheck.C) {
    helper := WrongSetUpSuiteArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.n, Equals, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpSuiteArgHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgHelper\\.SetUpSuite argument " +
        "should be \\*gocheck\\.C\n"

    c.Check(output.value, Matches, expected)
}


// -----------------------------------------------------------------------
// Nice errors also when tests or fixture have wrong arg count.

func (s *FixtureS) TestPanicOnWrongTestArgCount(c *gocheck.C) {
    helper := WrongTestArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.calls[0], Equals, "SetUpSuite")
    c.Check(helper.calls[1], Equals, "SetUpTest")
    c.Check(helper.calls[2], Equals, "TearDownTest")
    c.Check(helper.calls[3], Equals, "SetUpTest")
    c.Check(helper.calls[4], Equals, "Test2")
    c.Check(helper.calls[5], Equals, "TearDownTest")
    c.Check(helper.calls[6], Equals, "TearDownSuite")
    c.Check(helper.n, Equals, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:[0-9]+: " +
                "WrongTestArgCountHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgCountHelper\\.Test1 argument " +
                "should be \\*gocheck\\.C\n"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArgCount(c *gocheck.C) {
    helper := WrongSetUpTestArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.n, Equals, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpTestArgCountHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgCountHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.C\n"

    c.Check(output.value, Matches, expected)
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArgCount(c *gocheck.C) {
    helper := WrongSetUpSuiteArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.n, Equals, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpSuiteArgCountHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgCountHelper" +
        "\\.SetUpSuite argument should be \\*gocheck\\.C\n"

    c.Check(output.value, Matches, expected)
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

func (s *WrongTestArgCountHelper) Test1(c *gocheck.C, i int) {
}

type WrongSetUpTestArgCountHelper struct {
    FixtureHelper
}

func (s *WrongSetUpTestArgCountHelper) SetUpTest(c *gocheck.C, i int) {
}

type WrongSetUpSuiteArgCountHelper struct {
    FixtureHelper
}

func (s *WrongSetUpSuiteArgCountHelper) SetUpSuite(c *gocheck.C, i int) {
}


// -----------------------------------------------------------------------
// Ensure fixture doesn't run without tests.

type NoTestsHelper struct{
    hasRun bool
}

func (s *NoTestsHelper) SetUpSuite(c *gocheck.C) {
    s.hasRun = true
}

func (s *NoTestsHelper) TearDownSuite(c *gocheck.C) {
    s.hasRun = true
}

func (s *FixtureS) TestFixtureDoesntRunWithoutTests(c *gocheck.C) {
    helper := NoTestsHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Check(helper.hasRun, Equals, false)
}


// -----------------------------------------------------------------------
// Verify that checks and assertions work correctly inside the fixture.

type FixtureCheckHelper struct{
    fail string
    completed bool
}

func (s *FixtureCheckHelper) SetUpSuite(c *gocheck.C) {
    switch s.fail {
        case "SetUpSuiteAssert":
            c.Assert(false, Equals, true)
        case "SetUpSuiteCheck":
            c.Check(false, Equals, true)
    }
    s.completed = true
}

func (s *FixtureCheckHelper) SetUpTest(c *gocheck.C) {
    switch s.fail {
        case "SetUpTestAssert":
            c.Assert(false, Equals, true)
        case "SetUpTestCheck":
            c.Check(false, Equals, true)
    }
    s.completed = true
}

func (s *FixtureCheckHelper) Test(c *gocheck.C) {
    // Do nothing.
}

func (s *FixtureS) TestSetUpSuiteCheck(c *gocheck.C) {
    helper := FixtureCheckHelper{fail: "SetUpSuiteCheck"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Assert(output.value, Matches,
             "\n---+\n" +
             "FAIL: fixture_test\\.go:[0-9]+: " +
             "FixtureCheckHelper\\.SetUpSuite\n\n" +
             "fixture_test\\.go:[0-9]+:\n" +
             "\\.+ Check\\(obtained, Equals, expected\\):\n" +
             "\\.+ Obtained \\(bool\\): false\n" +
             "\\.+ Expected \\(bool\\): true\n\n")
    c.Assert(helper.completed, Equals, true)
}

func (s *FixtureS) TestSetUpSuiteAssert(c *gocheck.C) {
    helper := FixtureCheckHelper{fail: "SetUpSuiteAssert"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.Assert(output.value, Matches,
             "\n---+\n" +
             "FAIL: fixture_test\\.go:[0-9]+: " +
             "FixtureCheckHelper\\.SetUpSuite\n\n" +
             "fixture_test\\.go:[0-9]+:\n" +
             "\\.+ Assert\\(obtained, Equals, expected\\):\n" +
             "\\.+ Obtained \\(bool\\): false\n" +
             "\\.+ Expected \\(bool\\): true\n\n")
    c.Assert(helper.completed, Equals, false)
}
