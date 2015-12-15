package check_test

import (
	"fmt"
	"path/filepath"
	"runtime"

	. "gopkg.in/check.v1"
)

var _ = Suite(&checkReporterS{})

type checkReporterS struct {
	testFile string
}

func (s *checkReporterS) SetUpSuite(c *C) {
	_, fileName, _, ok := runtime.Caller(0)
	c.Assert(ok, Equals, true)
	s.testFile = filepath.Base(fileName)
}

func (s *checkReporterS) TestWrite(c *C) {
	testString := "test string"
	output := String{}

	dummyStream := true
	dummyVerbose := true
	r := NewCheckReporter(&output, dummyStream, dummyVerbose)

	r.Write([]byte(testString))
	c.Assert(output.value, Equals, testString)
}

func (s *checkReporterS) TestWriteCallStartedWithStreamFlag(c *C) {
	stream := true
	output := String{}

	dummyVerbose := true
	r := NewCheckReporter(&output, stream, dummyVerbose)

	r.WriteStarted(c)
	expected := fmt.Sprintf("START: %s:\\d+: %s\n", s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *checkReporterS) TestWriteCallStartedWithoutStreamFlag(c *C) {
	stream := false
	output := String{}

	dummyVerbose := true
	r := NewCheckReporter(&output, stream, dummyVerbose)

	r.WriteStarted(c)
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

func (s *checkReporterS) TestWriteCallProblemWithStreamFlag(c *C) {
	for _, testLabel := range problemTests {
		stream := true
		output := String{}

		dummyVerbose := true
		r := NewCheckReporter(&output, stream, dummyVerbose)

		writeProblem(testLabel, r, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\n\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *checkReporterS) TestWriteCallProblemWithoutStreamFlag(c *C) {
	for _, testLabel := range problemTests {
		stream := false
		output := String{}

		dummyVerbose := true
		r := NewCheckReporter(&output, stream, dummyVerbose)

		writeProblem(testLabel, r, c)
		expected := fmt.Sprintf(""+
			"\n"+
			"----------------------------------------------------------------------\n"+
			"%s: %s:\\d+: %s\n\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *checkReporterS) TestWriteCallProblemWithoutStreamFlagWithLog(c *C) {
	for _, testLabel := range problemTests {
		testLog := "test log"
		stream := false
		output := String{}

		dummyVerbose := true
		r := NewCheckReporter(&output, stream, dummyVerbose)

		c.Log(testLog)
		writeProblem(testLabel, r, c)
		expected := fmt.Sprintf(""+
			"\n"+
			"----------------------------------------------------------------------\n"+
			"%s: %s:\\d+: %s\n\n%s\n", testLabel, s.testFile, c.TestName(), testLog)
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

func (s *checkReporterS) TestWriteCallSuccessWithStreamFlag(c *C) {
	for _, testLabel := range successTests {
		stream := true
		output := String{}

		dummyVerbose := true
		r := NewCheckReporter(&output, stream, dummyVerbose)

		writeSuccess(testLabel, r, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *checkReporterS) TestWriteCallSuccessWithStreamFlagAndReason(c *C) {
	for _, testLabel := range successTests {
		testReason := "test skip reason"
		stream := true
		output := String{}

		dummyVerbose := true
		r := NewCheckReporter(&output, stream, dummyVerbose)
		c.FakeSkip(testReason)

		writeSuccess(testLabel, r, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s \\(%s\\)\t\\d\\.\\d+s\n\n",
			testLabel, s.testFile, c.TestName(), testReason)
		c.Check(output.value, Matches, expected)
	}
}

func (s *checkReporterS) TestWriteCallSuccessWithoutStreamFlagWithVerboseFlag(c *C) {
	for _, testLabel := range successTests {
		stream := false
		verbose := true
		output := String{}

		r := NewCheckReporter(&output, stream, verbose)

		writeSuccess(testLabel, r, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *checkReporterS) TestWriteCallSuccessWithoutStreamFlagWithoutVerboseFlag(c *C) {
	for _, testLabel := range successTests {
		stream := false
		verbose := false
		output := String{}

		r := NewCheckReporter(&output, stream, verbose)

		writeSuccess(testLabel, r, c)
		c.Check(output.value, Equals, "")
	}
}
