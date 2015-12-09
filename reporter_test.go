package check_test

import (
	"fmt"

	. "gopkg.in/check.v1"
)

const testFile = "reporter_test.go"

var _ = Suite(&ReporterS{})

type ReporterS struct{}

func (s *ReporterS) TestWrite(c *C) {
	testString := "test string"
	output := String{}

	dummyStream := true
	dummyVerbose := true
	o := NewOutputWriter(&output, dummyStream, dummyVerbose)

	o.Write([]byte(testString))
	c.Assert(output.value, Equals, testString)
}

func (s *ReporterS) TestWriteCallStartedWithStreamFlag(c *C) {
	testLabel := "test started label"
	stream := true
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteCallStarted(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\n", testLabel, testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallStartedWithoutStreamFlag(c *C) {
	stream := false
	output := String{}

	dummyLabel := "dummy"
	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteCallStarted(dummyLabel, c)
	c.Assert(output.value, Equals, "")
}

func (s *ReporterS) TestWriteCallProblemWithStreamFlag(c *C) {
	testLabel := "test problem label"
	stream := true
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteCallProblem(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\n\n", testLabel, testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallProblemWithoutStreamFlag(c *C) {
	testLabel := "test problem label"
	stream := false
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteCallProblem(testLabel, c)
	expected := fmt.Sprintf(""+
		"\n"+
		"----------------------------------------------------------------------\n"+
		"%s: %s:\\d+: %s\n\n", testLabel, testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallProblemWithoutStreamFlagWithLog(c *C) {
	testLabel := "test problem label"
	testLog := "test log"
	stream := false
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	c.Log(testLog)
	o.WriteCallProblem(testLabel, c)
	expected := fmt.Sprintf(""+
		"\n"+
		"----------------------------------------------------------------------\n"+
		"%s: %s:\\d+: %s\n\n%s\n", testLabel, testFile, c.TestName(), testLog)
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallSuccessWithStreamFlag(c *C) {
	testLabel := "test success label"
	stream := true
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteCallSuccess(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n\n", testLabel, testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallSuccessWithStreamFlagAndReason(c *C) {
	testLabel := "test success label"
	testReason := "test skip reason"
	stream := true
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)
	c.FakeSkip(testReason)

	o.WriteCallSuccess(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s \\(%s\\)\t\\d\\.\\d+s\n\n",
		testLabel, testFile, c.TestName(), testReason)
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallSuccessWithoutStreamFlagWithVerboseFlag(c *C) {
	testLabel := "test success label"
	stream := false
	verbose := true
	output := String{}

	o := NewOutputWriter(&output, stream, verbose)

	o.WriteCallSuccess(testLabel, c)
	expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n", testLabel, testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallSuccessWithoutStreamFlagWithoutVerboseFlag(c *C) {
	testLabel := "test success label"
	stream := false
	verbose := false
	output := String{}

	o := NewOutputWriter(&output, stream, verbose)

	o.WriteCallSuccess(testLabel, c)
	c.Assert(output.value, Equals, "")
}
