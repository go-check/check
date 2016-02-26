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

package check

import (
	"fmt"
	"strings"
)

type BootstrapS struct{}

var boostrapS = Suite(&BootstrapS{})

func (s *BootstrapS) TestCountSuite(c *C) {
	suitesRun += 1
}

func (s *BootstrapS) TestFailedAndFail(c *C) {
	if c.Failed() {
		critical("c.Failed() must be false first!")
	}
	c.Fail()
	if !c.Failed() {
		critical("c.Fail() didn't put the test in a failed state!")
	}
	c.Succeed()
}

func (s *BootstrapS) TestFailedAndSucceed(c *C) {
	c.Fail()
	c.Succeed()
	if c.Failed() {
		critical("c.Succeed() didn't put the test back in a non-failed state")
	}
}

func (s *BootstrapS) TestLogAndGetTestLog(c *C) {
	c.Log("Hello there!")
	log := c.GetTestLog()
	if log != "Hello there!\n" {
		critical(fmt.Sprintf("Log() or GetTestLog() is not working! Got: %#v", log))
	}
}

func (s *BootstrapS) TestLogfAndGetTestLog(c *C) {
	c.Logf("Hello %v", "there!")
	log := c.GetTestLog()
	if log != "Hello there!\n" {
		critical(fmt.Sprintf("Logf() or GetTestLog() is not working! Got: %#v", log))
	}
}

func (s *BootstrapS) TestRunShowsErrors(c *C) {
	output := String{}
	Run(&FailHelper{}, &RunConf{Output: &output})
	if strings.Index(output.value, "Expected failure!") == -1 {
		critical(fmt.Sprintf("RunWithWriter() output did not contain the "+
			"expected failure! Got: %#v",
			output.value))
	}
}

func (s *BootstrapS) TestRunDoesntShowSuccesses(c *C) {
	output := String{}
	Run(&SuccessHelper{}, &RunConf{Output: &output})
	if strings.Index(output.value, "Expected success!") != -1 {
		critical(fmt.Sprintf("RunWithWriter() output contained a successful "+
			"test! Got: %#v",
			output.value))
	}
}
