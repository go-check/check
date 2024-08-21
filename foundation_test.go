// These tests check that the foundations of gocheck are working properly.
// They already assume that fundamental failing is working already, though,
// since this was tested in bootstrap_test.go. Even then, some care may
// still have to be taken when using external functions, since they should
// of course not rely on functionality tested here.

package check_test

import (
	"log"
	"os"

	"gopkg.in/check.v1"
)

// -----------------------------------------------------------------------
// Foundation test suite.

type FoundationS struct{}

var foundationS = check.Suite(&FoundationS{})

func (s *FoundationS) TestCountSuite(c *check.C) {
	suitesRun += 1
}

// -----------------------------------------------------------------------
// Check minimum *log.Logger interface provided by *check.C.

type minLogger interface {
	Output(calldepth int, s string) error
}

func (s *FoundationS) TestMinLogger(c *check.C) {
	var logger minLogger
	logger = log.New(os.Stderr, "", 0)
	logger = c
	logger.Output(0, "Hello there")
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

var embeddedS = check.Suite(&EmbeddedS{})

func (s *EmbeddedS) TestCountSuite(c *check.C) {
	suitesRun += 1
}

func (s *EmbeddedInternalS) TestMethod(c *check.C) {
	c.Error("TestMethod() of the embedded type was called!?")
}

func (s *EmbeddedS) TestMethod(c *check.C) {
	// http://code.google.com/p/go/issues/detail?id=906
	c.Check(s.called, check.Equals, false) // Go issue 906 is affecting the runner?
	s.called = true
}
