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
                "PANIC: gocheck_test\\.go:FixtureHelper.Test1\n\n" +
                "\\.\\.\\. Panic: Test1 \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
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
                "PANIC: gocheck_test\\.go:FixtureHelper\\.SetUpTest\n\n" +
                "\\.\\.\\. Panic: SetUpTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.SetUpTest\n" +
                "\n-+\n" +
                "PANIC: gocheck_test\\.go:FixtureHelper\\.Test1\n\n" +
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
                "PANIC: gocheck_test\\.go:FixtureHelper.TearDownTest\n\n" +
                "\\.\\.\\. Panic: TearDownTest \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.TearDownTest\n" +
                "\n-+\n" +
                "PANIC: gocheck_test\\.go:FixtureHelper\\.Test1\n\n" +
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
                "PANIC: gocheck_test\\.go:FixtureHelper.SetUpSuite\n\n" +
                "\\.\\.\\. Panic: SetUpSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
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
                "PANIC: gocheck_test\\.go:FixtureHelper.TearDownSuite\n\n" +
                "\\.\\.\\. Panic: TearDownSuite \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
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
    t.CheckEqual(helper.fh.calls[0], "SetUpSuite")
    t.CheckEqual(helper.fh.calls[1], "SetUpTest")
    t.CheckEqual(helper.fh.calls[2], "TearDownTest")
    t.CheckEqual(helper.fh.calls[3], "SetUpTest")
    t.CheckEqual(helper.fh.calls[4], "Test2")
    t.CheckEqual(helper.fh.calls[5], "TearDownTest")
    t.CheckEqual(helper.fh.calls[6], "TearDownSuite")
    t.CheckEqual(helper.fh.n, 7)

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
    t.CheckEqual(helper.fh.n, 0)

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
    t.CheckEqual(helper.fh.n, 0)

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
    t.CheckEqual(helper.fh.calls[0], "SetUpSuite")
    t.CheckEqual(helper.fh.calls[1], "SetUpTest")
    t.CheckEqual(helper.fh.calls[2], "TearDownTest")
    t.CheckEqual(helper.fh.calls[3], "SetUpTest")
    t.CheckEqual(helper.fh.calls[4], "Test2")
    t.CheckEqual(helper.fh.calls[5], "TearDownTest")
    t.CheckEqual(helper.fh.calls[6], "TearDownSuite")
    t.CheckEqual(helper.fh.n, 7)

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
    t.CheckEqual(helper.fh.n, 0)

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
    t.CheckEqual(helper.fh.n, 0)

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
// Helper test suites with wrong function arguments.

// Wasn't for issue 906 in Go, we could embed FixtureHelper here
// rather than redefining all functions. :-(
type WrongTestArgHelper struct {
    fh FixtureHelper
}

func (s *WrongTestArgHelper) Test1(t int) {
}

func (s *WrongTestArgHelper) SetUpSuite(f *gocheck.F) {
    s.fh.SetUpSuite(f)
}

func (s *WrongTestArgHelper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *WrongTestArgHelper) SetUpTest(f *gocheck.F) {
    s.fh.SetUpTest(f)
}

func (s *WrongTestArgHelper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *WrongTestArgHelper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}


// ----

type WrongSetUpTestArgHelper struct {
    fh FixtureHelper
}

func (s *WrongSetUpTestArgHelper) SetUpTest(t int) {
}

func (s *WrongSetUpTestArgHelper) SetUpSuite(f *gocheck.F) {
    s.fh.SetUpSuite(f)
}

func (s *WrongSetUpTestArgHelper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *WrongSetUpTestArgHelper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *WrongSetUpTestArgHelper) Test1(t *gocheck.T) {
    s.fh.Test1(t)
}

func (s *WrongSetUpTestArgHelper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}


// ----

type WrongSetUpSuiteArgHelper struct {
    fh FixtureHelper
}

func (s *WrongSetUpSuiteArgHelper) SetUpSuite(t int) {
}

func (s *WrongSetUpSuiteArgHelper) SetUpTest(f *gocheck.F) {
    s.fh.SetUpTest(f)
}

func (s *WrongSetUpSuiteArgHelper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *WrongSetUpSuiteArgHelper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *WrongSetUpSuiteArgHelper) Test1(t *gocheck.T) {
    s.fh.Test1(t)
}

func (s *WrongSetUpSuiteArgHelper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}


// ----

type WrongTestArgCountHelper struct {
    fh FixtureHelper
}

func (s *WrongTestArgCountHelper) Test1(t *gocheck.T, i int) {
    s.fh.Test1(t)
}

func (s *WrongTestArgCountHelper) SetUpSuite(f *gocheck.F) {
    s.fh.SetUpSuite(f)
}

func (s *WrongTestArgCountHelper) SetUpTest(f *gocheck.F) {
    s.fh.SetUpTest(f)
}

func (s *WrongTestArgCountHelper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *WrongTestArgCountHelper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *WrongTestArgCountHelper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}


// ----

type WrongSetUpTestArgCountHelper struct {
    fh FixtureHelper
}

func (s *WrongSetUpTestArgCountHelper) SetUpTest(f *gocheck.F, i int) {
    s.fh.SetUpTest(f)
}

func (s *WrongSetUpTestArgCountHelper) SetUpSuite(f *gocheck.F) {
    s.fh.SetUpSuite(f)
}

func (s *WrongSetUpTestArgCountHelper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *WrongSetUpTestArgCountHelper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *WrongSetUpTestArgCountHelper) Test1(t *gocheck.T) {
    s.fh.Test1(t)
}

func (s *WrongSetUpTestArgCountHelper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}


// ----

type WrongSetUpSuiteArgCountHelper struct {
    fh FixtureHelper
}

func (s *WrongSetUpSuiteArgCountHelper) SetUpSuite(f *gocheck.F, i int) {
    s.fh.SetUpSuite(f)
}

func (s *WrongSetUpSuiteArgCountHelper) SetUpTest(f *gocheck.F) {
    s.fh.SetUpTest(f)
}

func (s *WrongSetUpSuiteArgCountHelper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *WrongSetUpSuiteArgCountHelper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *WrongSetUpSuiteArgCountHelper) Test1(t *gocheck.T) {
    s.fh.Test2(t)
}

func (s *WrongSetUpSuiteArgCountHelper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}


/*
type Helper struct {
    fh FixtureHelper
}

func (s *Helper) SetUpSuite(f *gocheck.F) {
    s.fh.SetUpSuite(f)
}

func (s *Helper) SetUpTest(f *gocheck.F) {
    s.fh.SetUpTest(f)
}

func (s *Helper) TearDownSuite(f *gocheck.F) {
    s.fh.TearDownSuite(f)
}

func (s *Helper) TearDownTest(f *gocheck.F) {
    s.fh.TearDownTest(f)
}

func (s *Helper) Test1(t *gocheck.T) {
    s.fh.Test1(t)
}

func (s *Helper) Test2(t *gocheck.T) {
    s.fh.Test2(t)
}
*/
