package check

func PrintLine(filename string, line int) (string, error) {
	return printLine(filename, line)
}

func Indent(s, with string) string {
	return indent(s, with)
}

func (c *C) FakeSkip(reason string) {
	c.reason = reason
}
