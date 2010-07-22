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


// -----------------------------------------------------------------------
// Basic fixture ordering verification.

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


// -----------------------------------------------------------------------
// Check the behavior when panics occur within tests and fixtures.

func (s *FixtureS) TestPanicOnTest(t *gocheck.T) {
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
                "PANIC: fixture_test\\.go:FixtureHelper.Test1\n\n" +
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

func (s *FixtureS) TestPanicOnSetUpTest(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "SetUpTest"}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "TearDownTest")
    t.CheckEqual(helper.calls[3], "TearDownSuite")
    t.CheckEqual(helper.n, 4)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:FixtureHelper\\.SetUpTest\n\n" +
                "\\.\\.\\. Panic: SetUpTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpTest\n" +
                "\n-+\n" +
                "PANIC: fixture_test\\.go:FixtureHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: Fixture has panicked " +
                "\\(see related PANIC\\)\n$"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnTearDownTest(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "TearDownTest"}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "Test1")
    t.CheckEqual(helper.calls[3], "TearDownTest")
    t.CheckEqual(helper.calls[4], "TearDownSuite")
    t.CheckEqual(helper.n, 5)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:FixtureHelper.TearDownTest\n\n" +
                "\\.\\.\\. Panic: TearDownTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownTest\n" +
                "\n-+\n" +
                "PANIC: fixture_test\\.go:FixtureHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: Fixture has panicked " +
                "\\(see related PANIC\\)\n$"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnSetUpSuite(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "SetUpSuite"}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "TearDownSuite")
    t.CheckEqual(helper.n, 2)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:FixtureHelper.SetUpSuite\n\n" +
                "\\.\\.\\. Panic: SetUpSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpSuite\n$"

    // XXX Changing the expression above to not match breaks Go. WTF?

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnTearDownSuite(t *gocheck.T) {
    helper := FixtureHelper{panicOn: "TearDownSuite"}
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
                "PANIC: fixture_test\\.go:FixtureHelper.TearDownSuite\n\n" +
                "\\.\\.\\. Panic: TearDownSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*fixture_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownSuite\n$"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression:", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// A wrong argument on a test or fixture will produce a nice error.

func (s *FixtureS) TestPanicOnWrongTestArg(t *gocheck.T) {
    helper := WrongTestArgHelper{}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "TearDownTest")
    t.CheckEqual(helper.calls[3], "SetUpTest")
    t.CheckEqual(helper.calls[4], "Test2")
    t.CheckEqual(helper.calls[5], "TearDownTest")
    t.CheckEqual(helper.calls[6], "TearDownSuite")
    t.CheckEqual(helper.n, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:WrongTestArgHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgHelper\\.Test1 argument " +
                "should be \\*gocheck\\.T\n"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArg(t *gocheck.T) {
    helper := WrongSetUpTestArgHelper{}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:WrongSetUpTestArgHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.F\n"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArg(t *gocheck.T) {
    helper := WrongSetUpSuiteArgHelper{}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:WrongSetUpSuiteArgHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgHelper\\.SetUpSuite argument " +
        "should be \\*gocheck\\.F\n"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// Nice errors also when tests or fixture have wrong arg count.

func (s *FixtureS) TestPanicOnWrongTestArgCount(t *gocheck.T) {
    helper := WrongTestArgCountHelper{}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.calls[0], "SetUpSuite")
    t.CheckEqual(helper.calls[1], "SetUpTest")
    t.CheckEqual(helper.calls[2], "TearDownTest")
    t.CheckEqual(helper.calls[3], "SetUpTest")
    t.CheckEqual(helper.calls[4], "Test2")
    t.CheckEqual(helper.calls[5], "TearDownTest")
    t.CheckEqual(helper.calls[6], "TearDownSuite")
    t.CheckEqual(helper.n, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:WrongTestArgCountHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgCountHelper\\.Test1 argument " +
                "should be \\*gocheck\\.T\n"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArgCount(t *gocheck.T) {
    helper := WrongSetUpTestArgCountHelper{}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:WrongSetUpTestArgCountHelper" +
        "\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgCountHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.F\n"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArgCount(t *gocheck.T) {
    helper := WrongSetUpSuiteArgCountHelper{}
    output := String{}
    gocheck.RunWithWriter(&helper, &output)
    t.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:WrongSetUpSuiteArgCountHelper" +
        "\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgCountHelper" +
        "\\.SetUpSuite argument should be \\*gocheck\\.F\n"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
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

func (s *FixtureHelper) SetUpSuite(f *gocheck.F) {
    s.trace("SetUpSuite")
}

func (s *FixtureHelper) TearDownSuite(f *gocheck.F) {
    s.trace("TearDownSuite")
}

func (s *FixtureHelper) SetUpTest(f *gocheck.F) {
    s.trace("SetUpTest")
}

func (s *FixtureHelper) TearDownTest(f *gocheck.F) {
    s.trace("TearDownTest")
}

func (s *FixtureHelper) Test1(t *gocheck.T) {
    s.trace("Test1")
}

func (s *FixtureHelper) Test2(t *gocheck.T) {
    s.trace("Test2")
}


// -----------------------------------------------------------------------
// Helper test suites with wrong function arguments.

type WrongTestArgHelper struct {
    FixtureHelper // XXX Using this will explode due to issue 906 in Go.
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

func (s *WrongTestArgCountHelper) Test1(t *gocheck.T, i int) {
}

type WrongSetUpTestArgCountHelper struct {
    FixtureHelper
}

func (s *WrongSetUpTestArgCountHelper) SetUpTest(f *gocheck.F, i int) {
}

type WrongSetUpSuiteArgCountHelper struct {
    FixtureHelper
}

func (s *WrongSetUpSuiteArgCountHelper) SetUpSuite(f *gocheck.F, i int) {
}
