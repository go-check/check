package check

import "io"

type TestReporter interface {
	testReporter
}

func PrintLine(filename string, line int) (string, error) {
	return printLine(filename, line)
}

func Indent(s, with string) string {
	return indent(s, with)
}

func NewOutputWriter(writer io.Writer, verbosity uint8) *outputWriter {
	return newOutputWriter(writer, verbosity)
}

func (c *C) FakeSkip(reason string) {
	c.reason = reason
}
