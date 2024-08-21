// Integration tests

package check_test

import (
	. "gopkg.in/check.v1"
)

// -----------------------------------------------------------------------
// Integration test suite.

type integrationS struct{}

var _ = Suite(&integrationS{})

type integrationTestHelper struct{}

func (s *integrationTestHelper) TestMultiLineStringEqualFails(c *C) {
	c.Check("foo\nbar\nbaz\nboom\n", Equals, "foo\nbaar\nbaz\nboom\n")
}

func (s *integrationTestHelper) TestStringEqualFails(c *C) {
	c.Check("foo", Equals, "bar")
}

func (s *integrationTestHelper) TestIntEqualFails(c *C) {
	c.Check(42, Equals, 43)
}

type complexStruct struct {
	r, i int
}

func (s *integrationTestHelper) TestStructEqualFails(c *C) {
	c.Check(complexStruct{1, 2}, Equals, complexStruct{3, 4})
}

func (s *integrationS) TestCountSuite(c *C) {
	suitesRun += 1
}

func (s *integrationS) TestOutput(c *C) {
	exitCode, output := runHelperSuite(c, "integrationTestHelper")
	c.Check(exitCode, Equals, 1)

	c.Check(output.Status("integrationTestHelper/TestIntEqualFails"), Equals, "fail")
	c.Check(output.Logs("integrationTestHelper/TestIntEqualFails"), Equals, ""+
		"    integration_test.go:27: \n"+
		"            c.Check(42, Equals, 43)\n"+
		"        ... obtained int = 42\n"+
		"        ... expected int = 43\n")

	c.Check(output.Status("integrationTestHelper/TestMultiLineStringEqualFails"), Equals, "fail")
	c.Check(output.Logs("integrationTestHelper/TestMultiLineStringEqualFails"), Equals, ""+
		"    integration_test.go:19: \n"+
		"            c.Check(\"foo\\nbar\\nbaz\\nboom\\n\", Equals, \"foo\\nbaar\\nbaz\\nboom\\n\")\n"+
		"        ... obtained string = \"\" +\n"+
		"        ...     \"foo\\n\" +\n"+
		"        ...     \"bar\\n\" +\n"+
		"        ...     \"baz\\n\" +\n"+
		"        ...     \"boom\\n\"\n"+
		"        ... expected string = \"\" +\n"+
		"        ...     \"foo\\n\" +\n"+
		"        ...     \"baar\\n\" +\n"+
		"        ...     \"baz\\n\" +\n"+
		"        ...     \"boom\\n\"\n"+
		"        ... String difference:\n"+
		"        ...     [1]: \"bar\" != \"baar\"\n")

	c.Check(output.Status("integrationTestHelper/TestStringEqualFails"), Equals, "fail")
	c.Check(output.Logs("integrationTestHelper/TestStringEqualFails"), Equals, ""+
		"    integration_test.go:23: \n"+
		"            c.Check(\"foo\", Equals, \"bar\")\n"+
		"        ... obtained string = \"foo\"\n"+
		"        ... expected string = \"bar\"\n")

	c.Check(output.Status("integrationTestHelper/TestStructEqualFails"), Equals, "fail")
	c.Check(output.Logs("integrationTestHelper/TestStructEqualFails"), Equals, ""+
		"    integration_test.go:35: \n"+
		"            c.Check(complexStruct{1, 2}, Equals, complexStruct{3, 4})\n"+
		"        ... obtained check_test.complexStruct = check_test.complexStruct{r:1, i:2}\n"+
		"        ... expected check_test.complexStruct = check_test.complexStruct{r:3, i:4}\n"+
		"        ... Difference:\n"+
		"        ...     r: 1 != 3\n"+
		"        ...     i: 2 != 4\n")
}
