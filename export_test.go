package check

import "io"

type Reporter interface {
	reporter
}

func PrintLine(filename string, line int) (string, error) {
	return printLine(filename, line)
}

func Indent(s, with string) string {
	return indent(s, with)
}

func NewCheckReporter(writer io.Writer, stream, verbose bool) *checkReporter {
	return newCheckReporter(writer, stream, verbose)
}

func (c *C) FakeSkip(reason string) {
	c.reason = reason
}
