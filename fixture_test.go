/*
Gocheck - A rich testing framework for Go

Copyright (c) 2010, Gustavo Niemeyer <gustavo@niemeyer.net>

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of the copyright holder nor the names of its
      contributors may be used to endorse or promote products derived from
      this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Tests for the behavior of the test fixture system.

package gocheck_test


import (
    "gocheck"
    "regexp"
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
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "SetUpTest")
    c.CheckEqual(helper.calls[5], "Test2")
    c.CheckEqual(helper.calls[6], "TearDownTest")
    c.CheckEqual(helper.calls[7], "TearDownSuite")
    c.CheckEqual(helper.n, 8)
}


// -----------------------------------------------------------------------
// Check the behavior when panics occur within tests and fixtures.

func (s *FixtureS) TestPanicOnTest(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "Test1"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "SetUpTest")
    c.CheckEqual(helper.calls[5], "Test2")
    c.CheckEqual(helper.calls[6], "TearDownTest")
    c.CheckEqual(helper.calls[7], "TearDownSuite")
    c.CheckEqual(helper.n, 8)

    expected := "^\n-+\n" +
                "PANIC: gocheck_test\\.go:[0-9]+: FixtureHelper.Test1\n\n" +
                "\\.\\.\\. Panic: Test1 \\(PC=[xA-F0-9]+\\)\n\n" +
                ".+:[0-9]+\n" +
                "  in runtime.panic\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.trace\n" +
                ".*gocheck_test.go:[0-9]+\n" +
                "  in FixtureHelper.Test1\n$"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression:", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnSetUpTest(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "SetUpTest"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "TearDownTest")
    c.CheckEqual(helper.calls[3], "TearDownSuite")
    c.CheckEqual(helper.n, 4)

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

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression:", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnTearDownTest(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "TearDownTest"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "TearDownSuite")
    c.CheckEqual(helper.n, 5)

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

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression:", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnSetUpSuite(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "SetUpSuite"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "TearDownSuite")
    c.CheckEqual(helper.n, 2)

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

    // XXX Changing the expression above to not match breaks Go. WTF?

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression:", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnTearDownSuite(c *gocheck.C) {
    helper := FixtureHelper{panicOn: "TearDownSuite"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "Test1")
    c.CheckEqual(helper.calls[3], "TearDownTest")
    c.CheckEqual(helper.calls[4], "SetUpTest")
    c.CheckEqual(helper.calls[5], "Test2")
    c.CheckEqual(helper.calls[6], "TearDownTest")
    c.CheckEqual(helper.calls[7], "TearDownSuite")
    c.CheckEqual(helper.n, 8)

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

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression:", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// A wrong argument on a test or fixture will produce a nice error.

func (s *FixtureS) TestPanicOnWrongTestArg(c *gocheck.C) {
    helper := WrongTestArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "TearDownTest")
    c.CheckEqual(helper.calls[3], "SetUpTest")
    c.CheckEqual(helper.calls[4], "Test2")
    c.CheckEqual(helper.calls[5], "TearDownTest")
    c.CheckEqual(helper.calls[6], "TearDownSuite")
    c.CheckEqual(helper.n, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:[0-9]+: " +
                "WrongTestArgHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgHelper\\.Test1 argument " +
                "should be \\*gocheck\\.C\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression: ", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArg(c *gocheck.C) {
    helper := WrongSetUpTestArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpTestArgHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.C\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression: ", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArg(c *gocheck.C) {
    helper := WrongSetUpSuiteArgHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpSuiteArgHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgHelper\\.SetUpSuite argument " +
        "should be \\*gocheck\\.C\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression: ", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// Nice errors also when tests or fixture have wrong arg count.

func (s *FixtureS) TestPanicOnWrongTestArgCount(c *gocheck.C) {
    helper := WrongTestArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.calls[0], "SetUpSuite")
    c.CheckEqual(helper.calls[1], "SetUpTest")
    c.CheckEqual(helper.calls[2], "TearDownTest")
    c.CheckEqual(helper.calls[3], "SetUpTest")
    c.CheckEqual(helper.calls[4], "Test2")
    c.CheckEqual(helper.calls[5], "TearDownTest")
    c.CheckEqual(helper.calls[6], "TearDownSuite")
    c.CheckEqual(helper.n, 7)

    expected := "^\n-+\n" +
                "PANIC: fixture_test\\.go:[0-9]+: " +
                "WrongTestArgCountHelper\\.Test1\n\n" +
                "\\.\\.\\. Panic: WrongTestArgCountHelper\\.Test1 argument " +
                "should be \\*gocheck\\.C\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression: ", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpTestArgCount(c *gocheck.C) {
    helper := WrongSetUpTestArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpTestArgCountHelper\\.SetUpTest\n\n" +
        "\\.\\.\\. Panic: WrongSetUpTestArgCountHelper\\.SetUpTest argument " +
        "should be \\*gocheck\\.C\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression: ", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}

func (s *FixtureS) TestPanicOnWrongSetUpSuiteArgCount(c *gocheck.C) {
    helper := WrongSetUpSuiteArgCountHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.CheckEqual(helper.n, 0)

    expected :=
        "^\n-+\n" +
        "PANIC: fixture_test\\.go:[0-9]+: " +
        "WrongSetUpSuiteArgCountHelper\\.SetUpSuite\n\n" +
        "\\.\\.\\. Panic: WrongSetUpSuiteArgCountHelper" +
        "\\.SetUpSuite argument should be \\*gocheck\\.C\n"

    matched, err := regexp.MatchString(expected, output.value)
    if err != nil {
        c.Error("Bad expression: ", expected)
    } else if !matched {
        c.Error("Panic not logged properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// Helper test suites with wrong function arguments.

type WrongTestArgHelper struct {
    FixtureHelper
}

func (s *WrongTestArgHelper) Test1(t int) {
}

// ----

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
    c.CheckEqual(helper.hasRun, false)
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
            c.AssertEqual(false, true)
        case "SetUpSuiteCheck":
            c.CheckEqual(false, true)
    }
    s.completed = true
}

func (s *FixtureCheckHelper) SetUpTest(c *gocheck.C) {
    switch s.fail {
        case "SetUpTestAssert":
            c.AssertEqual(false, true)
        case "SetUpTestCheck":
            c.CheckEqual(false, true)
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
    c.AssertMatch(output.value,
                  "\n---+\n" +
                  "FAIL: fixture_test\\.go:[0-9]+: " +
                  "FixtureCheckHelper\\.SetUpSuite\n\n" +
                  "fixture_test\\.go:[0-9]+:\n" +
                  "\\.+ CheckEqual\\(obtained, expected\\):\n" +
                  "\\.+ Obtained \\(bool\\): false\n" +
                  "\\.+ Expected \\(bool\\): true\n\n")
    c.AssertEqual(helper.completed, true)
}

func (s *FixtureS) TestSetUpSuiteAssert(c *gocheck.C) {
    helper := FixtureCheckHelper{fail: "SetUpSuiteAssert"}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.AssertMatch(output.value,
                  "\n---+\n" +
                  "FAIL: fixture_test\\.go:[0-9]+: " +
                  "FixtureCheckHelper\\.SetUpSuite\n\n" +
                  "fixture_test\\.go:[0-9]+:\n" +
                  "\\.+ AssertEqual\\(obtained, expected\\):\n" +
                  "\\.+ Obtained \\(bool\\): false\n" +
                  "\\.+ Expected \\(bool\\): true\n\n")
    c.AssertEqual(helper.completed, false)
}
