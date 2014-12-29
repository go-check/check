package check

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"
)

// TODO: start test suite
type outputWriter interface {
	Write(content []byte) (n int, err error)
	WriteCallStarted(label string, c *C)
	WriteCallProblem(label string, c *C)
	WriteCallSuccess(label string, c *C)
	StreamEnabled() bool
}

/*************** Plain writer *****************/

type plainWriter struct {
	outputWriter
	m                    sync.Mutex
	writer               io.Writer
	wroteCallProblemLast bool
	stream               bool
	verbose              bool
}

func newPlainWriter(writer io.Writer, stream, verbose bool) *plainWriter {
	return &plainWriter{writer: writer, stream: stream, verbose: verbose}
}

func (w *plainWriter) StreamEnabled() bool { return w.stream }

func (w *plainWriter) Write(content []byte) (n int, err error) {
	w.m.Lock()
	n, err = w.writer.Write(content)
	w.m.Unlock()
	return
}

func (w *plainWriter) WriteCallStarted(label string, c *C) {
	if w.stream {
		header := renderCallHeader(label, c, "", "\n")
		w.m.Lock()
		w.writer.Write([]byte(header))
		w.m.Unlock()
	}
}

func (w *plainWriter) WriteCallProblem(label string, c *C) {
	var prefix string
	if !w.stream {
		prefix = "\n-----------------------------------" +
			"-----------------------------------\n"
	}
	header := renderCallHeader(label, c, prefix, "\n\n")
	w.m.Lock()
	w.wroteCallProblemLast = true
	w.writer.Write([]byte(header))
	if !w.stream {
		c.logb.WriteTo(w.writer)
	}
	w.m.Unlock()
}

func (w *plainWriter) WriteCallSuccess(label string, c *C) {
	if w.stream || (w.verbose && c.kind == testKd) {
		// TODO Use a buffer here.
		var suffix string
		if c.reason != "" {
			suffix = " (" + c.reason + ")"
		}
		if c.status == succeededSt {
			suffix += "\t" + c.timerString()
		}
		suffix += "\n"
		if w.stream {
			suffix += "\n"
		}
		header := renderCallHeader(label, c, "", suffix)
		w.m.Lock()
		// Resist temptation of using line as prefix above due to race.
		if !w.stream && w.wroteCallProblemLast {
			header = "\n-----------------------------------" +
				"-----------------------------------\n" +
				header
		}
		w.wroteCallProblemLast = false
		w.writer.Write([]byte(header))
		w.m.Unlock()
	}
}

func renderCallHeader(label string, c *C, prefix, suffix string) string {
	pc := c.method.PC()
	return fmt.Sprintf("%s%s: %s: %s%s", prefix, label, niceFuncPath(pc),
		niceFuncName(pc), suffix)
}

/*************** xUnit writer *****************/
// TODO: Write mrthod can collect data for std-out
type xunitReport struct {
	suites []xunitSuite `xml:"testsuites>testsuite,omitempty"`
}

type xunitSuite struct {
	Package   string    `xml:"package,attr,omitempty"`
	Name      string    `xml:"name,attr,omitempty"`
	Classname string    `xml:"classname,attr,omitempty"`
	Time      float64   `xml:"time,attr"`
	Timestamp time.Time `xml:"timestamp,attr"`

	Tests    uint64 `xml:"tests,attr"`
	Failures uint64 `xml:"failures,attr"`
	Errors   uint64 `xml:"errors,attr"`
	Skipped  uint64 `xml:"skipped,attr"`

	Properties []xunitSuiteProperty `xml:"properties>property"` //TODO: test
	Testcases  []xunitTestcase      `xml:"testcase"`

	SystemOut string `xml:"system-out,omitempty"`
	SystemErr string `xml:"system-err,omitempty"`
}

type xunitSuiteProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type xunitTestcase struct {
	Name      string               `xml:"name,attr,omitempty"`
	Classname string               `xml:"classname,attr,omitempty"`
	Time      float64              `xml:"time,attr"`
	Failure   *xunitTestcaseResult `xml:"failure,omitempty"`
	Error     *xunitTestcaseResult `xml:"error,omitempty"`
}

type xunitTestcaseResult struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Value   string `xml:",chardata"`
}

type xunitWriter struct {
	outputWriter
	m       sync.Mutex
	writer  io.Writer
	stream  bool
	verbose bool

	systemOut io.Writer
}

func newXunitWriter(writer io.Writer, stream, verbose bool) *plainWriter {
	return &xunitWriter{
		writer:    writer,
		systemOut: bytes.Buffer{},
		stream:    stream,
		verbose:   verbose,
	}
}

func (w *xunitWriter) GetReport() string {
}

func (w *xunitWriter) Write(content []byte) (n int, err error) {
	w.m.Lock()
	n, err = w.systemOut.Write(content)
	w.m.Unlock()
	return
}
func (w *xunitWriter) WriteCallStarted(label string, c *C) {}
func (w *xunitWriter) WriteCallProblem(label string, c *C) {}
func (w *xunitWriter) WriteCallSuccess(label string, c *C) {}
func (w *xunitWriter) StreamEnabled() bool                 {}
