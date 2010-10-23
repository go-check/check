// These initial tests are for bootstrapping.  They verify that we can
// basically use the testing infrastructure itself to check if the test
// system is working.
// 
// These tests use will break down the test runner badly in case of
// errors because if they simply fail, we can't be sure the developer
// will ever see anything (because failing means the failing system
// somehow isn't working! :-)
//
// Do not assume *any* internal functionality works as expected besides
// what's actually tested here.

package gocheck_test


import (
    "gocheck"
    "strings"
    "fmt"
)


type BootstrapS struct{}

var boostrapS = gocheck.Suite(&BootstrapS{})

func (s *BootstrapS) TestCountSuite(t *gocheck.T) {
    suitesRun += 1
}

func (s *BootstrapS) TestFailedAndFail(t *gocheck.T) {
    if t.Failed() {
        critical("t.Failed() must be false first!")
    }
    t.Fail()
    if !t.Failed() {
        critical("t.Fail() didn't put the test in a failed state!")
    }
    t.Succeed()
}

func (s *BootstrapS) TestFailedAndSucceed(t *gocheck.T) {
    t.Fail()
    t.Succeed()
    if t.Failed() {
        critical("t.Succeed() didn't put the test back in a non-failed state")
    }
}

func (s *BootstrapS) TestLogAndGetLog(t *gocheck.T) {
    t.Log("Hello there!")
    log := t.GetLog()
    if log != "Hello there!\n" {
        critical(fmt.Sprintf("Log() or GetLog() is not working! Got: %#v", log))
    }
}

func (s *BootstrapS) TestLogfAndGetLog(t *gocheck.T) {
    t.Logf("Hello %v", "there!")
    log := t.GetLog()
    if log != "Hello there!\n" {
        critical(fmt.Sprintf("Logf() or GetLog() is not working! Got: %#v", log))
    }
}

func (s *BootstrapS) TestRunShowsErrors(t *gocheck.T) {
    output := String{}
    gocheck.Run(&FailHelper{}, &gocheck.RunConf{Output: &output})
    if strings.Index(output.value, "Expected failure!") == -1 {
        critical(fmt.Sprintf("RunWithWriter() output did not contain the " +
                             "expected failure! Got: %#v", output.value))
    }
}

func (s *BootstrapS) TestRunDoesntShowSuccesses(t *gocheck.T) {
    output := String{}
    gocheck.Run(&SuccessHelper{}, &gocheck.RunConf{Output: &output})
    if strings.Index(output.value, "Expected success!") != -1 {
        critical(fmt.Sprintf("RunWithWriter() output contained a successful " +
                             "test! Got: %#v", output.value))
    }
}
