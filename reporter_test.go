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

func (s *reporterS) TestStartTestWithHighVerbosity(c *C) {
	var verbosity uint8 = 2
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.StartTest(c)
	expected := fmt.Sprintf("START: %s:\\d+: %s\n", s.testFile, c.TestName())
	c.Assert(output.value, Matches, expected)
}

func (s *reporterS) TestStartTestLowVerbosity(c *C) {
	var verbosity uint8 = 1
	output := String{}

	o := NewOutputWriter(&output, verbosity)

	o.StartTest(c)
	c.Assert(output.value, Equals, "")
}

var problemTests = []string{"FAIL", "PANIC"}

func addProblem(label string, r TestReporter, c *C) {
	if label == "FAIL" {
		r.AddFailure(c)
	} else if label == "PANIC" {
		r.AddError(c)
	} else {
		panic("Unknown problem: " + label)
	}
}

func (s *reporterS) TestAddProblemWitHighVerbosity(c *C) {
	for _, testLabel := range problemTests {
		var verbosity uint8 = 2
		output := String{}
		
		o := NewOutputWriter(&output, verbosity)

		addProblem(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\n\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *reporterS) TestAddProblemWithLowVerbosity(c *C) {
	for _, testLabel := range problemTests {
		var verbosity uint8 = 1
		output := String{}

		o := NewOutputWriter(&output, verbosity)
		
		addProblem(testLabel, o, c)
		expected := fmt.Sprintf(""+
			"\n"+
			"----------------------------------------------------------------------\n"+
			"%s: %s:\\d+: %s\n\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *reporterS) TestAddProblemWithLowVerbosityWithLog(c *C) {
	for _, testLabel := range problemTests {
		testLog := "test log"
		var verbosity uint8 = 1
		output := String{}

		o := NewOutputWriter(&output, verbosity)

		c.Log(testLog)
		addProblem(testLabel, o, c)

		expected := fmt.Sprintf(""+
			"\n"+
			"----------------------------------------------------------------------\n"+
			"%s: %s:\\d+: %s\n\n%s\n", testLabel, s.testFile, c.TestName(), testLog)
		c.Check(output.value, Matches, expected)
	}
}

var successTests = []string{"PASS", "SKIP", "FAIL EXPECTED", "MISS"}

func addSuccess(label string, r TestReporter, c *C) {
	if label == "PASS" {
		r.AddSuccess(c)
	} else if label == "SKIP" {
		r.AddSkip(c)
	} else if label == "FAIL EXPECTED" {
		r.AddExpectedFailure(c)
	} else if label == "MISS" {
		r.AddMissed(c)
	} else {
		panic("Unknown success: " + label)
	}
}

func (s *reporterS) TestAddSuccessWithHighVerbosity(c *C) {
	for _, testLabel := range successTests {
		var verbosity uint8 = 2
		output := String{}

		o := NewOutputWriter(&output, verbosity)

		addSuccess(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *reporterS) TestAddSuccessWithHighVerbosityAndReason(c *C) {
	for _, testLabel := range successTests {
		testReason := "test skip reason"
		var verbosity uint8 = 2
		output := String{}

		o := NewOutputWriter(&output, verbosity)
		c.FakeSkip(testReason)

		addSuccess(testLabel, o, c)
		expected := fmt.Sprintf("%s: %s:\\d+: %s \\(%s\\)\t\\d\\.\\d+s\n\n",
			testLabel, s.testFile, c.TestName(), testReason)
		c.Check(output.value, Matches, expected)
	}
}

func (s *reporterS) TestAddSuccessWithLowVerbosity(c *C) {
	for _, testLabel := range successTests {
		var verbosity uint8 = 1
		output := String{}

		o := NewOutputWriter(&output, verbosity)

		addSuccess(testLabel, o, c)

		expected := fmt.Sprintf("%s: %s:\\d+: %s\t\\d\\.\\d+s\n", testLabel, s.testFile, c.TestName())
		c.Check(output.value, Matches, expected)
	}
}

func (s *reporterS) TestAddSuccessWithoutVerbosity(c *C) {
	for _, testLabel := range successTests {
		var verbosity uint8 = 0
		output := String{}

		o := NewOutputWriter(&output, verbosity)

		addSuccess(testLabel, o, c)
		c.Check(output.value, Equals, "")
	}
}
