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
	stream := true
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteStarted(c)
	expected := fmt.Sprintf("START: %s:\\d+: %s\n", testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *ReporterS) TestWriteCallStartedWithoutStreamFlag(c *C) {
	stream := false
	output := String{}

	dummyVerbose := true
	o := NewOutputWriter(&output, stream, dummyVerbose)

	o.WriteStarted(c)
	c.Assert(output.value, Equals, "")
}

var problemTests = []string{"FAIL", "PANIC"}

func writeProblem(label string, r Reporter, c *C) {
	if label == "FAIL" {
		r.WriteFailure(c)
	} else if label == "PANIC" {
		r.WriteError(c)
	} else {
		panic("Unknown problem: " + label)
	}
}

func (s *ReporterS) TestWriteCallProblemWithStreamFlag(c *C) {
	for _, testLabel := range problemTests {
		stream := true
		output := String{}

		dummyVerbose := true
		o := NewOutputWriter(&output, stream, dummyVerbose)

		writeProblem(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\n\n", testLabel, testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *ReporterS) TestWriteCallProblemWithoutStreamFlag(c *C) {
	for _, testLabel := range problemTests {
		stream := false
		output := String{}

		dummyVerbose := true
		o := NewOutputWriter(&output, stream, dummyVerbose)

		writeProblem(testLabel, o, c)
		expected := fmt.Sprintf(""+
			"\n"+
			"----------------------------------------------------------------------\n"+
			"%s: %s:\\d+: %s\n\n", testLabel, testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *ReporterS) TestWriteCallProblemWithoutStreamFlagWithLog(c *C) {
	for _, testLabel := range problemTests {
		testLog := "test log"
		stream := false
		output := String{}

		dummyVerbose := true
		o := NewOutputWriter(&output, stream, dummyVerbose)

		c.Log(testLog)
		writeProblem(testLabel, o, c)
		expected := fmt.Sprintf(""+
			"\n"+
			"----------------------------------------------------------------------\n"+
			"%s: %s:\\d+: %s\n\n%s\n", testLabel, testFile, c.TestName(), testLog)
		c.Check(output.value, Matches, expected)
	}
}

var successTests = []string{"PASS", "SKIP", "FAIL EXPECTED", "MISS"}

func writeSuccess(label string, r Reporter, c *C) {
	if label == "PASS" {
		r.WriteSuccess(c)
	} else if label == "SKIP" {
		r.WriteSkip(c)
	} else if label == "FAIL EXPECTED" {
		r.WriteExpectedFailure(c)
	} else if label == "MISS" {
		r.WriteMissed(c)
	} else {
		panic("Unknown success: " + label)
	}
}

func (s *ReporterS) TestWriteCallSuccessWithStreamFlag(c *C) {
	for _, testLabel := range successTests {
		stream := true
		output := String{}

		dummyVerbose := true
		o := NewOutputWriter(&output, stream, dummyVerbose)

		writeSuccess(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n\n", testLabel, testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *ReporterS) TestWriteCallSuccessWithStreamFlagAndReason(c *C) {
	for _, testLabel := range successTests {
		testReason := "test skip reason"
		stream := true
		output := String{}

		dummyVerbose := true
		o := NewOutputWriter(&output, stream, dummyVerbose)
		c.FakeSkip(testReason)

		writeSuccess(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s \\(%s\\)\t\\d\\.\\d+s\n\n",
			testLabel, testFile, c.TestName(), testReason)
		c.Check(output.value, Matches, expected)
	}
}

func (s *ReporterS) TestWriteCallSuccessWithoutStreamFlagWithVerboseFlag(c *C) {
	for _, testLabel := range successTests {
		stream := false
		verbose := true
		output := String{}

		o := NewOutputWriter(&output, stream, verbose)

		writeSuccess(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n", testLabel, testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *ReporterS) TestWriteCallSuccessWithoutStreamFlagWithoutVerboseFlag(c *C) {
	for _, testLabel := range successTests {
		stream := false
		verbose := false
		output := String{}

		o := NewOutputWriter(&output, stream, verbose)

		writeSuccess(testLabel, o, c)
		c.Check(output.value, Equals, "")
	}
}
