package check

import (
	"fmt"
	"io"
	"sync"
)

type reporter interface {
	io.Writer
	WriteStarted(*C)
	WriteFailure(*C)
	WriteError(*C)
	WriteSuccess(*C)
	WriteSkip(*C)
	WriteExpectedFailure(*C)
	WriteMissed(*C)
	Stream() bool
}

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

type outputWriter struct {
	m                    sync.Mutex
	writer               io.Writer
	wroteCallProblemLast bool
	stream               bool
	verbose              bool
}

func newOutputWriter(writer io.Writer, stream, verbose bool) *outputWriter {
	return &outputWriter{writer: writer, stream: stream, verbose: verbose}
}

func (ow *outputWriter) Stream() bool {
	return ow.stream
}

func (ow *outputWriter) Write(content []byte) (n int, err error) {
	ow.m.Lock()
	n, err = ow.writer.Write(content)
	ow.m.Unlock()
	return
}

func (ow *outputWriter) WriteStarted(c *C) {
	if ow.Stream() {
		header := renderCallHeader("START", c, "", "\n")
		ow.m.Lock()
		ow.writer.Write([]byte(header))
		ow.m.Unlock()
	}
}

func (ow *outputWriter) WriteFailure(c *C) {
	ow.writeProblem("FAIL", c)
}

func (ow *outputWriter) WriteError(c *C) {
	ow.writeProblem("PANIC", c)
}

func (ow *outputWriter) writeProblem(label string, c *C) {
	var prefix string
	if !ow.Stream() {
		prefix = "\n-----------------------------------" +
			"-----------------------------------\n"
	}
	header := renderCallHeader(label, c, prefix, "\n\n")
	ow.m.Lock()
	ow.wroteCallProblemLast = true
	ow.writer.Write([]byte(header))
	if !ow.Stream() {
		c.logb.WriteTo(ow.writer)
	}
	ow.m.Unlock()
}

func (ow *outputWriter) WriteSuccess(c *C) {
	ow.writeSuccess("PASS", c)
}

func (ow *outputWriter) WriteSkip(c *C) {
	ow.writeSuccess("SKIP", c)
}

func (ow *outputWriter) WriteExpectedFailure(c *C) {
	ow.writeSuccess("FAIL EXPECTED", c)
}

func (ow *outputWriter) WriteMissed(c *C) {
	ow.writeSuccess("MISS", c)
}

func (ow *outputWriter) writeSuccess(label string, c *C) {
	if ow.Stream() || (ow.verbose && c.kind == testKd) {
		// TODO Use a buffer here.
		var suffix string
		if c.reason != "" {
			suffix = " (" + c.reason + ")"
		}
		if c.status() == succeededSt {
			suffix += "\t" + c.timerString()
		}
		suffix += "\n"
		if ow.Stream() {
			suffix += "\n"
		}
		header := renderCallHeader(label, c, "", suffix)
		ow.m.Lock()
		// Resist temptation of using line as prefix above due to race.
		if !ow.Stream() && ow.wroteCallProblemLast {
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
