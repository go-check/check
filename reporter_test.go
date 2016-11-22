package check_test

import (
	"fmt"
	"path/filepath"
	"runtime"

	. "github.com/elopio/check"
)

var _ = Suite(&reporterS{})

type reporterS struct {
	testFile string
}

func (s *reporterS) SetUpSuite(c *C) {
	_, fileName, _, ok := runtime.Caller(0)
	c.Assert(ok, Equals, true)
	s.testFile = filepath.Base(fileName)
}

func (s *reporterS) TestWriteCallStartedWithStreamFlag(c *C) {
	testLabel := "test started label"
	var verbosity uint8 = 2
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.WriteCallStarted(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\n", testLabel, s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallStartedWithoutStreamFlag(c *C) {
	var verbosity uint8 = 1
	output := String{}

	dummyLabel := "dummy"
	o := NewOutputWriter(&output, verbosity)

	o.WriteCallStarted(dummyLabel, c)
	c.Assert(output.value, Equals, "")
}

func (s *reporterS) TestWriteCallProblemWithStreamFlag(c *C) {
	testLabel := "test problem label"
	var verbosity uint8 = 2
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.WriteCallProblem(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\n\n", testLabel, s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallProblemWithoutStreamFlag(c *C) {
	testLabel := "test problem label"
	var verbosity uint8 = 1
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.WriteCallProblem(testLabel, c)
	expected := fmt.Sprintf(""+
		"\n"+
		"----------------------------------------------------------------------\n"+
		"%s: %s:\\d+: %s\n\n", testLabel, s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallProblemWithoutStreamFlagWithLog(c *C) {
	testLabel := "test problem label"
	testLog := "test log"
	var verbosity uint8 = 1
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	c.Log(testLog)
	o.WriteCallProblem(testLabel, c)
	expected := fmt.Sprintf(""+
		"\n"+
		"----------------------------------------------------------------------\n"+
		"%s: %s:\\d+: %s\n\n%s\n", testLabel, s.testFile, c.TestName(), testLog)
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallSuccessWithStreamFlag(c *C) {
	testLabel := "test success label"
	var verbosity uint8 = 2
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.WriteCallSuccess(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n\n", testLabel, s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallSuccessWithStreamFlagAndReason(c *C) {
	testLabel := "test success label"
	testReason := "test skip reason"
	var verbosity uint8 = 2
	output := String{}

	o := NewOutputWriter(&output, verbosity)
	c.FakeSkip(testReason)

	o.WriteCallSuccess(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s \\(%s\\)\t\\d\\.\\d+s\n\n",
		testLabel, s.testFile, c.TestName(), testReason)
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallSuccessWithoutStreamFlagWithVerboseFlag(c *C) {
	testLabel := "test success label"
	var verbosity uint8 = 1
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.WriteCallSuccess(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n", testLabel, s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestWriteCallSuccessWithoutStreamFlagWithoutVerboseFlag(c *C) {
	testLabel := "test success label"
	var verbosity uint8 = 0
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.WriteCallSuccess(testLabel, c)
	c.Assert(output.value, Equals, "")
}
