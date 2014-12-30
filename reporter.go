package check

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"sync"
	"time"
)

type reporter interface {
	GetReport() (string, error)
}

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

	m sync.Mutex
}

func (s *xunitSuite) TestFail(tc xunitTestcase, message, value string) {
	tc.Failure = &xunitTestcaseResult{
		Message: message,
		Value:   value,
	}

	s.m.Lock()
	s.Tests++
	s.Failures++
	s.Testcases = append(s.Testcases, tc)
	s.m.Unlock()
}
func (s *xunitSuite) TestError(tc xunitTestcase, message, value string) {
	tc.Error = &xunitTestcaseResult{
		Message: message,
		Value:   value,
	}

	s.m.Lock()
	s.Tests++
	s.Errors++
	s.Testcases = append(s.Testcases, tc)
	s.m.Unlock()
}

func (s *xunitSuite) TestSkip(tc xunitTestcase) {
	s.m.Lock()
	s.Tests++
	s.Skipped++
	s.Testcases = append(s.Testcases, tc)
	s.m.Unlock()
}

func (s *xunitSuite) TestSuccess(tc xunitTestcase) {
	s.m.Lock()
	s.Tests++
	s.Testcases = append(s.Testcases, tc)
	s.m.Unlock()
}

type xunitSuiteProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type xunitTestcase struct {
	Name      string  `xml:"name,attr,omitempty"`
	Classname string  `xml:"classname,attr,omitempty"`
	Time      float64 `xml:"time,attr"`

	Failure *xunitTestcaseResult `xml:"failure,omitempty"`
	Error   *xunitTestcaseResult `xml:"error,omitempty"`
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
	suites  map[string]*xunitSuite

	systemOut io.Writer
}

func newXunitWriter(writer io.Writer, stream, verbose bool) *xunitWriter {
	return &xunitWriter{
		writer:    writer,
		systemOut: &bytes.Buffer{},
		stream:    stream,
		verbose:   verbose,
	}
}

func (w *xunitWriter) GetReport() (out string, err error) {
	report := xunitReport{}
	report.suites = make([]xunitSuite, len(w.suites), 0)
	for k := range w.suites {
		report.suites = append(report.suites, *w.suites[k])
	}

	var buf []byte
	buf, err = xml.Marshal(report)
	if err != nil {
		out = string(buf)
	}

	return
}

func (w *xunitWriter) Write(content []byte) (n int, err error) {
	w.m.Lock()
	n, err = w.systemOut.Write(content)
	w.m.Unlock()
	return
}

func (w *xunitWriter) WriteCallStarted(label string, c *C) {
	w.getSuite(c) // init suite if not yet existing
}

func (w *xunitWriter) WriteTestSkipped(label string, c *C) {
	res := w.newTestcase(c)
	w.getSuite(c).TestSkip(res)
}

func (w *xunitWriter) WriteCallFailure(label string, c *C) {
	res := w.newTestcase(c)
	w.getSuite(c).TestFail(res, label, c.logb.String())
}

func (w *xunitWriter) WriteCallError(label string, c *C) {
	res := w.newTestcase(c)
	w.getSuite(c).TestError(res, label, c.logb.String())
}

func (w *xunitWriter) WriteCallSuccess(label string, c *C) {
	res := w.newTestcase(c)
	w.getSuite(c).TestSuccess(res)
}

func (w *xunitWriter) StreamEnabled() bool { return w.stream }

func (w *xunitWriter) getSuite(c *C) (suite *xunitSuite) {
	var ok bool
	suiteName := c.method.suiteName()
	w.m.Lock()
	if suite, ok = w.suites[suiteName]; !ok {
		suite = &xunitSuite{
			Name:      suiteName,
			Timestamp: c.startTime,
		}
		w.suites[suiteName] = suite
	}
	w.m.Unlock()

	return
}

func (w *xunitWriter) newTestcase(c *C) xunitTestcase {
	return xunitTestcase{
		Name:      c.testName,
		Classname: c.method.suiteName(),
		Time:      time.Since(c.startTime).Seconds(),
	}
}
