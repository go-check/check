// These tests verify the inner workings of the helper methods associated
// with gocheck.T.

package gocheck_test

import (
    "gocheck"
    "testing"
    "strings"
    "fmt"
)


func TestHelpers(t *testing.T) {
    gocheck.RunTestingT(&HelpersS{}, t)
}


// -----------------------------------------------------------------------
// Helpers test suite.

type HelpersS struct{}

func (s *HelpersS) TestErrorf(t *gocheck.T) {
    // Do not use checkState() here.  It depends on Errorf() working.
    t.Errorf("Error %v!", "message")
    failed := t.Failed()
    t.Succeed()
    if log := t.GetLog(); log != "Error message!\n" {
        t.Logf("Errorf() hasn't logged the message. Got: %#v", log)
        t.Fail()
    }
    if !failed {
        t.Logf("Errorf() didn't put the test in a failed state")
        t.Fail()
    }
}

func (s *HelpersS) TestError(t *gocheck.T) {
    // XXX Should Error() and Errorf() include the file/line? Probably!
    t.Error("Error", "message!")
    checkState(t, nil,
               &expectedState{
                    name: "Error(`Error`, `message!`)",
                    failed: true,
                    log: "Error message!\n",
               })
}

func (s *HelpersS) TestCheckEqualSucceeding(t *gocheck.T) {
    result := t.CheckEqual(10, 10)
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 10)",
                    result: true,
               })
}

func (s *HelpersS) TestCheckEqualFailing(t *gocheck.T) {
    log := fmt.Sprintf(
        "\n... %d:CheckEqual(A, B): A != B\n" +
        "... A: 10\n" +
        "... B: 20\n\n",
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

func (s *HelpersS) TestCheckEqualWithMessage(t *gocheck.T) {
    log := fmt.Sprintf(
        "\n... %d:CheckEqual(A, B): A != B\n" +
        "... A: 10\n" +
        "... B: 20\n" +
        "... That's clearly.. WRONG!\n\n",
        getMyLine()+1)
    result := t.CheckEqual(10, 20, "That's clearly.. ", "WRONG!")
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 20)",
                    result: false,
                    failed: true,
                    log: log,
               })
}

func (s *HelpersS) TestCheckEqualFailingInDifferentFile(t *gocheck.T) {
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

func (s *HelpersS) TestFailureHeader(t *gocheck.T) {
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


// -----------------------------------------------------------------------
// Helper which checks the state of the test and ensures that it matches
// the given expectations.  Depends on t.Errorf() working, so shouldn't
// be used to test this one function.

type expectedState struct {
    name string
    result interface{}
    failed bool
    log string
}

// Verify the state of the test.  Note that since this also verifies if
// the test is supposed to be in a failed state, no other checks should
// be done in addition to what is being tested.
func checkState(t *gocheck.T, result interface{}, expected *expectedState) {
    failed := t.Failed()
    t.Succeed()
    if log := t.GetLog(); log != expected.log {
        t.Errorf("%s logged %#v rather than %#v",
                 expected.name, log, expected.log)
    }
    if result != expected.result {
        t.Errorf("%s returned %#v rather than %#v",
                 expected.name, result, expected.result)
    }
    if failed != expected.failed {
        if failed {
            t.Errorf("%s has failed when it shouldn't", expected.name)
        } else {
            t.Errorf("%s has not failed when it should", expected.name)
        }
    }
}
