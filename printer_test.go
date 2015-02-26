package check_test

import (
	. "github.com/masukomi/check"
)

var _ = Suite(&PrinterS{})

type PrinterS struct{}

func (s *PrinterS) TestCountSuite(c *C) {
	suitesRun += 1
}

var printTestFuncLine int

func init() {
	printTestFuncLine = getMyLine() + 3
}

func printTestFunc() {
	println(1)  // Comment1
	if 2 == 2 { // Comment2
		println(3) // Comment3
	}
	switch 5 {
	case 6:
		println(6) // Comment6
	}
	println(7)
	switch interface{}(9).(type) { // Comment9
	case int:
		println(10)
		println(11)
	}
	select {
	case <-(chan bool)(nil):
		println(14)
		println(15)
	default:
		println(16)
		println(17)
	}
	println(19,
		20)
	_ = func() {
		println(21)
		println(22)
	}
	println(24, func() {
		println(25)
	})
	// Leading comment
	// with multiple lines.
	println(29) // Comment29
}

var printLineTests = []struct {
	line   int
	output string
}{
	{1, "println(1) // Comment1"},
	{2, "if 2 == 2 { // Comment2\n    ...\n}"},
	{3, "println(3) // Comment3"},
	{5, "switch 5 {\n...\n}"},
	{6, "case 6:\n    ... // Comment6"},
	{9, "println(7)"},
	{10, "switch interface{}(9).(type) { // Comment9\n...\n}"},
	{11, "case int:\n    ...\n    println(11)"},
	{16, "case <-(chan bool)(nil):\n    ...\n    println(15)"},
	{18, "println(15)"},
	{19, "default:\n    ...\n    println(17)"},
	{21, "println(17)"},
	{23, "println(19,\n    20)"},
	{24, "println(19,\n    20)"},
	{25, "_ = func() {\n    println(21)\n    println(22)\n}"},
	{27, "println(22)"},
	{29, "println(24, func() {\n    println(25)\n})"},
	{30, "println(25)"},
	{31, "println(24, func() {\n    println(25)\n})"},
	{34, "// Leading comment\n// with multiple lines.\nprintln(29) // Comment29"},
}

func (s *PrinterS) TestPrintLine(c *C) {
	for _, test := range printLineTests {
		output, err := PrintLine("printer_test.go", printTestFuncLine+test.line)
		c.Assert(err, IsNil)
		c.Assert(output, Equals, test.output)
	}
}

var indentTests = []struct {
	in, out string
}{
	{"", ""},
	{"\n", "\n"},
	{"a", ">>>a"},
	{"a\n", ">>>a\n"},
	{"a\nb", ">>>a\n>>>b"},
	{" ", ">>> "},
}

func (s *PrinterS) TestIndent(c *C) {
	for _, test := range indentTests {
		out := Indent(test.in, ">>>")
		c.Assert(out, Equals, test.out)
	}

}
