// These tests check that the foundations of gocheck are working properly.
// They already assume that fundamental failing is working already, though,
// since this was tested in bootstrap_test.go. Even then, some care may
// still have to be taken when using external functions, since they should
// of course not rely on functionality tested here.

package gocheck_test

import (
    "gocheck"
    "testing"
    "strings"
    "fmt"
)


func TestFoundation(t *testing.T) {
    gocheck.RunTestingT(&FoundationS{}, t)
}


// -----------------------------------------------------------------------
// Foundation test suite.

type FoundationS struct{}


func (s *FoundationS) TestErrorf(t *gocheck.T) {
    // Do not use checkState() here.  It depends on Errorf() working.
    expectedLog := fmt.Sprintf("... %d:Error: Error message!\n", getMyLine()+1)
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
    expectedLog := fmt.Sprintf("... %d:Error: Error message!\n", getMyLine()+1)
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
    t.Log("This shouldn't be logged!")
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
    t.Log("This shouldn't be logged!")
}

func (s *FoundationS) TestFailureHeader(t *gocheck.T) {
    output := String{}
    failHelper := FailHelper{}
    gocheck.RunWithWriter(&failHelper, &output)
    header := fmt.Sprintf(
        "-----------------------------------" +
        "-----------------------------------\n" +
        "FAIL: gocheck_test.go:TestLogAndFail\n")
        //"FAIL: gocheck_test.go:%d:TestLogAndFail\n", failHelper.testLine)
        // How to find the first line of a function?
    if strings.Index(output.value, header) == -1 {
        t.Errorf("Failure didn't print a proper header.\n" +
                 "... Got:\n%s... Expected something with:\n%s",
                 output.value, header)
    }
}

func (s *FoundationS) TestCallerLoggingInDifferentFile(t *gocheck.T) {
    result, line := checkEqualWrapper(t, 10, 20)
    log := fmt.Sprintf(
        "\n... gocheck_test.go:%d:CheckEqual(A, B): A != B\n" +
        "... A: 10\n" +
        "... B: 20\n\n",
        line)
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 20)",
                    result: false,
                    failed: true,
                    log: log,
               })
}
