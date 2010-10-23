// These tests check that the foundations of gocheck are working properly.
// They already assume that fundamental failing is working already, though,
// since this was tested in bootstrap_test.go. Even then, some care may
// still have to be taken when using external functions, since they should
// of course not rely on functionality tested here.

package gocheck_test

import (
    "testing"
    "gocheck"
    "strings"
    "fmt"
)


// -----------------------------------------------------------------------
// Foundation test suite.

type FoundationS struct{}

var foundationS = gocheck.Suite(&FoundationS{})

func (s *FoundationS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}

func (s *FoundationS) TestErrorf(t *gocheck.T) {
    // Do not use checkState() here.  It depends on Errorf() working.
    expectedLog := fmt.Sprintf("foundation_test.go:%d:\n" +
                               "... Error: Error message!\n", getMyLine()+1)
    t.Errorf("Error %v!", "message")
    failed := t.Failed()
    t.Succeed()
    if log := t.GetLog(); log != expectedLog {
        t.Logf("Errorf() logged %#v rather than %#v", log, expectedLog)
        t.Fail()
    }
    if !failed {
        t.Logf("Errorf() didn't put the test in a failed state")
        t.Fail()
    }
}

func (s *FoundationS) TestError(t *gocheck.T) {
    expectedLog := fmt.Sprintf("foundation_test.go:%d:\n" +
                               "... Error: Error message!\n", getMyLine()+1)
    t.Error("Error ", "message!")
    checkState(t, nil,
               &expectedState{
                    name: "Error(`Error `, `message!`)",
                    failed: true,
                    log: expectedLog,
               })
}

func (s *FoundationS) TestFailNow(t *gocheck.T) {
    defer (func() {
        if !t.Failed() {
            t.Error("FailNow() didn't fail the test")
        } else {
            t.Succeed()
            t.CheckEqual(t.GetLog(), "")
        }
    })()

    t.FailNow()
    t.Log("FailNow() didn't stop the test")
}

func (s *FoundationS) TestSucceedNow(t *gocheck.T) {
    defer (func() {
        if t.Failed() {
            t.Error("SucceedNow() didn't succeed the test")
        }
        t.CheckEqual(t.GetLog(), "")
    })()

    t.Fail()
    t.SucceedNow()
    t.Log("SucceedNow() didn't stop the test")
}

func (s *FoundationS) TestFailureHeader(t *gocheck.T) {
    output := String{}
    failHelper := FailHelper{}
    gocheck.Run(&failHelper, &gocheck.RunConf{Output: &output})
    header := fmt.Sprintf(
        "\n-----------------------------------" +
        "-----------------------------------\n" +
        "FAIL: gocheck_test.go:%d: FailHelper.TestLogAndFail\n",
        failHelper.testLine)
    if strings.Index(output.value, header) == -1 {
        t.Errorf("Failure didn't print a proper header.\n" +
                 "... Got:\n%s... Expected something with:\n%s",
                 output.value, header)
    }
}

func (s *FoundationS) TestFatal(t *gocheck.T) {
    var line int
    defer (func() {
        if !t.Failed() {
            t.Error("Fatal() didn't fail the test")
        } else {
            t.Succeed()
            t.CheckEqual(t.GetLog(),
                         fmt.Sprintf("foundation_test.go:%d:\n" +
                                     "... Error: Die now!\n", line))
        }
    })()

    line = getMyLine()+1
    t.Fatal("Die ", "now!")
    t.Log("Fatal() didn't stop the test")
}

func (s *FoundationS) TestFatalf(t *gocheck.T) {
    var line int
    defer (func() {
        if !t.Failed() {
            t.Error("Fatalf() didn't fail the test")
        } else {
            t.Succeed()
            t.CheckEqual(t.GetLog(),
                         fmt.Sprintf("foundation_test.go:%d:\n" +
                                     "... Error: Die now!\n", line))
        }
    })()

    line = getMyLine()+1
    t.Fatalf("Die %s!", "now")
    t.Log("Fatalf() didn't stop the test")
}


func (s *FoundationS) TestCallerLoggingInsideTest(t *gocheck.T) {
    log := fmt.Sprintf(
        "foundation_test.go:%d:\n" +
        "\\.\\.\\. CheckEqual\\(obtained, expected\\):\n" +
        "\\.\\.\\. Obtained \\(int\\): 10\n" +
        "\\.\\.\\. Expected \\(int\\): 20\n\n",
        getMyLine()+1)
    result := t.CheckEqual(10, 20)
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 20)",
                    result: false,
                    failed: true,
                    log: log,
               })
}

func (s *FoundationS) TestCallerLoggingInDifferentFile(t *gocheck.T) {
    result, line := checkEqualWrapper(t, 10, 20)
    testLine := getMyLine()-1
    log := fmt.Sprintf(
        "foundation_test.go:%d > gocheck_test.go:%d:\n" +
        "\\.\\.\\. CheckEqual\\(obtained, expected\\):\n" +
        "\\.\\.\\. Obtained \\(int\\): 10\n" +
        "\\.\\.\\. Expected \\(int\\): 20\n\n",
        testLine, line)
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 20)",
                    result: false,
                    failed: true,
                    log: log,
               })
}

// -----------------------------------------------------------------------
// ExpectFailure() inverts the logic of failure.

type ExpectFailureHelper struct{}

func (s *ExpectFailureHelper) TestFail(t *gocheck.T) {
    t.ExpectFailure("It booms!")
    t.Error("Boom!")
}

func (s *ExpectFailureHelper) TestSucceed(t *gocheck.T) {
    t.ExpectFailure("Bug #XYZ")
}

func (s *FoundationS) TestExpectFailure(t *gocheck.T) {
    helper := ExpectFailureHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})

    expected :=
        "^\n-+\n" +
        "FAIL: foundation_test\\.go:[0-9]+:" +
        " ExpectFailureHelper\\.TestSucceed\n\n" +
        "\\.\\.\\. Error: Test succeeded, but was expected to fail\n" +
        "\\.\\.\\. Reason: Bug #XYZ\n$"

    matched, err := testing.MatchString(expected, output.value)
    if err != "" {
        t.Error("Bad expression: ", expected)
    } else if !matched {
        t.Error("ExpectFailure() didn't log properly:\n", output.value)
    }
}


// -----------------------------------------------------------------------
// Ensure that suites with embedded types are working fine, including the
// the workaround for issue 906.

type EmbeddedInternalS struct {
    called bool
}

type EmbeddedS struct {
    EmbeddedInternalS
}

var embeddedS = gocheck.Suite(&EmbeddedS{})

func (s *EmbeddedS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}

func (s *EmbeddedInternalS) TestMethod(t *gocheck.T) {
    t.Error("TestMethod() of the embedded type was called!?")
}

func (s *EmbeddedS) TestMethod(t *gocheck.T) {
    // http://code.google.com/p/go/issues/detail?id=906
    t.CheckEqual(s.called, false,
                 "The bug described in issue 906 is affecting the runner")
    s.called = true
}
