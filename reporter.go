package check

import (
	"fmt"
	"io"
	"sync"
)

type testReporter interface {
	StartTest(*C)
	StopTest(*C)
	AddFailure(*C)
	AddError(*C)
	AddUnexpectedSuccess(*C)
	AddSuccess(*C)
	AddExpectedFailure(*C)
	AddSkip(*C)
	AddMissed(*C)
}

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

type outputWriter struct {
	m                    sync.Mutex
	writer               io.Writer
	wroteCallProblemLast bool
	verbosity            uint8
}

func newOutputWriter(writer io.Writer, verbosity uint8) *outputWriter {
	return &outputWriter{writer: writer, verbosity: verbosity}
}

func (ow *outputWriter) StartTest(c *C) {
	if ow.verbosity > 1 {
		header := renderCallHeader("START", c, "", "\n")
		ow.m.Lock()
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	}
}

func (ow *outputWriter) StopTest(c *C) {
}

func (ow *outputWriter) AddFailure(c *C) {
	ow.writeProblem("FAIL", c)
}

func (ow *outputWriter) AddError(c *C) {
	ow.writeProblem("PANIC", c)
}

func (ow *outputWriter) writeProblem(label string, c *C) {
	var prefix string
	if ow.verbosity < 2 {
		prefix = "\n-----------------------------------" +
			"-----------------------------------\n"
	}
	header := renderCallHeader(label, c, prefix, "\n\n")
	ow.m.Lock()
	ow.wroteCallProblemLast = true
	ow.writer.Write([]byte(header))
	if ow.verbosity < 2 {
		c.logb.WriteTo(ow.writer)
	}
	ow.m.Unlock()
}

func (ow *outputWriter) AddUnexpectedSuccess(c *C) {
}

func (ow *outputWriter) AddSuccess(c *C) {
	ow.writeSuccess("PASS", c)
}

func (ow *outputWriter) AddExpectedFailure(c *C) {
	ow.writeSuccess("FAIL EXPECTED", c)
}

func (ow *outputWriter) AddSkip(c *C) {
	ow.writeSuccess("SKIP", c)
}

func (ow *outputWriter) AddMissed(c *C) {
	ow.writeSuccess("MISS", c)
}
	
func (ow *outputWriter) writeSuccess(label string, c *C) {
	if ow.verbosity > 1 || (ow.verbosity == 1 && c.kind == testKd) {
		// TODO Use a buffer here.
		var suffix string
		if c.reason != "" {
			suffix = " (" + c.reason + ")"
		}
		if c.status() == succeededSt {
			suffix += "\t" + c.timerString()
		}
		suffix += "\n"
		if ow.verbosity > 1 {
			suffix += "\n"
		}
		header := renderCallHeader(label, c, "", suffix)
		ow.m.Lock()
		// Resist temptation of using line as prefix above due to race.
		if ow.verbosity < 2 && ow.wroteCallProblemLast {
			header = "\n-----------------------------------" +
				"-----------------------------------\n" +
				header
		}
		ow.wroteCallProblemLast = false
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	}
}

func renderCallHeader(label string, c *C, prefix, suffix string) string {
	pc := c.method.PC()
	return fmt.Sprintf("%s%s: %s: %s%s", prefix, label, niceFuncPath(pc),
		niceFuncName(pc), suffix)
}
