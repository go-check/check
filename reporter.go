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
// Reporters manage atomic output writing according to settings.

type checkReporter struct {
	m                    sync.Mutex
	writer               io.Writer
	wroteCallProblemLast bool
	stream               bool
	verbose              bool
}

func newCheckReporter(writer io.Writer, stream, verbose bool) *checkReporter {
	return &checkReporter{writer: writer, stream: stream, verbose: verbose}
}

func (r *checkReporter) Stream() bool {
	return r.stream
}

func (r *checkReporter) Write(content []byte) (n int, err error) {
	r.m.Lock()
	n, err = r.writer.Write(content)
	r.m.Unlock()
	return
}

func (r *checkReporter) WriteStarted(c *C) {
	if r.Stream() {
		header := renderCallHeader("START", c, "", "\n")
		r.m.Lock()
		r.writer.Write([]byte(header))
		r.m.Unlock()
	}
}

func (r *checkReporter) WriteFailure(c *C) {
	r.writeProblem("FAIL", c)
}

func (r *checkReporter) WriteError(c *C) {
	r.writeProblem("PANIC", c)
}

func (r *checkReporter) writeProblem(label string, c *C) {
	var prefix string
	if !r.Stream() {
		prefix = "\n-----------------------------------" +
			"-----------------------------------\n"
	}
	header := renderCallHeader(label, c, prefix, "\n\n")
	r.m.Lock()
	r.wroteCallProblemLast = true
	r.writer.Write([]byte(header))
	if !r.Stream() {
		c.logb.WriteTo(r.writer)
	}
	r.m.Unlock()
}

func (r *checkReporter) WriteSuccess(c *C) {
	r.writeSuccess("PASS", c)
}

func (r *checkReporter) WriteSkip(c *C) {
	r.writeSuccess("SKIP", c)
}

func (r *checkReporter) WriteExpectedFailure(c *C) {
	r.writeSuccess("FAIL EXPECTED", c)
}

func (r *checkReporter) WriteMissed(c *C) {
	r.writeSuccess("MISS", c)
}

func (r *checkReporter) writeSuccess(label string, c *C) {
	if r.Stream() || (r.verbose && c.kind == testKd) {
		// TODO Use a buffer here.
		var suffix string
		if c.reason != "" {
			suffix = " (" + c.reason + ")"
		}
		if c.status() == succeededSt {
			suffix += "\t" + c.timerString()
		}
		suffix += "\n"
		if r.Stream() {
			suffix += "\n"
		}
		header := renderCallHeader(label, c, "", suffix)
		r.m.Lock()
		// Resist temptation of using line as prefix above due to race.
		if !r.Stream() && r.wroteCallProblemLast {
			header = "\n-----------------------------------" +
				"-----------------------------------\n" +
				header
		}
		r.wroteCallProblemLast = false
		r.writer.Write([]byte(header))
		r.m.Unlock()
	}
}

func renderCallHeader(label string, c *C, prefix, suffix string) string {
	pc := c.method.PC()
	return fmt.Sprintf("%s%s: %s: %s%s", prefix, label, niceFuncPath(pc),
		niceFuncName(pc), suffix)
}
