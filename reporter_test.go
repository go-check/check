package check

/*************** xUnit writer tests *****************/
type XUnitTestSuite struct {
	writer *xunitWriter
}

var _ = Suite(&XUnitTestSuite{})

func (s *XUnitTestSuite) SetUpTest(c *C) {
	s.writer = newXunitWriter(nil, false)
}

func (s *XUnitTestSuite) TestSuccess(c *C) {
	s.writer.WriteCallSuccess("PASS", c)
	report, err := s.writer.GetReport()
	c.Assert(err, IsNil)

	match := "<testsuites>\n" +
		" +<testsuite .*name=\"XUnitTestSuite\" .*tests=\"1\" failures=\"0\" errors=\"0\" skipped=\"0\">\n" +
		" +<testcase name=\"XUnitTestSuite\\.TestSuccess\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*</testcase>\n" +
		" +</testsuite>\n" +
		"</testsuites>"

	c.Assert(string(report), Matches, match)
}

func (s *XUnitTestSuite) TestSkip(c *C) {
	s.writer.WriteCallSkipped("SKIP", c)
	report, err := s.writer.GetReport()
	c.Assert(err, IsNil)

	match := "<testsuites>\n" +
		" +<testsuite .*name=\"XUnitTestSuite\" .*tests=\"1\" failures=\"0\" errors=\"0\" skipped=\"1\">\n" +
		" +<testcase name=\"XUnitTestSuite\\.TestSkip\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<skipped>true</skipped>\n" +
		" +</testcase>\n" +
		" +</testsuite>\n" +
		"</testsuites>"

	c.Assert(string(report), Matches, match)
}

func (s *XUnitTestSuite) TestFail(c *C) {
	s.writer.WriteCallFailure("FAIL", c)
	report, err := s.writer.GetReport()
	c.Assert(err, IsNil)

	match := "<testsuites>\n" +
		" +<testsuite .*name=\"XUnitTestSuite\" .*tests=\"1\" failures=\"1\" errors=\"0\" skipped=\"0\">\n" +
		" +<testcase name=\"XUnitTestSuite\\.TestFail\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<failure message=\"FAIL\" type=\"go.failure\"></failure>\n" +
		" +</testcase>\n" +
		" +</testsuite>\n" +
		"</testsuites>"

	c.Assert(string(report), Matches, match)
}

func (s *XUnitTestSuite) TestError(c *C) {
	s.writer.WriteCallError("ERR", c)
	report, err := s.writer.GetReport()
	c.Assert(err, IsNil)

	match := "<testsuites>\n" +
		" +<testsuite .*name=\"XUnitTestSuite\" .*tests=\"1\" failures=\"0\" errors=\"1\" skipped=\"0\">\n" +
		" +<testcase name=\"XUnitTestSuite\\.TestError\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<error message=\"ERR\" type=\"go.error\"></error>\n" +
		" +</testcase>\n" +
		" +</testsuite>\n" +
		"</testsuites>"

	c.Assert(string(report), Matches, match)

}

func (s *XUnitTestSuite) TestCombine(c *C) {
	s.writer.WriteCallError("ERR", c)
	s.writer.WriteCallFailure("FAIL", c)
	s.writer.WriteCallSuccess("PASS", c)
	s.writer.WriteCallSuccess("PASS", c)
	s.writer.WriteCallFailure("FAIL", c)
	s.writer.WriteCallSkipped("SKIP", c)

	report, err := s.writer.GetReport()
	c.Assert(err, IsNil)

	match := "<testsuites>\n" +
		" +<testsuite .*name=\"XUnitTestSuite\" .*tests=\"6\" failures=\"2\" errors=\"1\" skipped=\"1\">\n" +
		" +<testcase name=\"XUnitTestSuite\\.TestCombine\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<error message=\"ERR\" type=\"go.error\"></error>\n" +
		" +</testcase>\n" +

		" +<testcase name=\"XUnitTestSuite\\.TestCombine\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<failure message=\"FAIL\" type=\"go.failure\"></failure>\n" +
		" +</testcase>\n" +

		" +<testcase name=\"XUnitTestSuite\\.TestCombine\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*</testcase>\n" +

		" +<testcase name=\"XUnitTestSuite\\.TestCombine\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*</testcase>\n" +

		" +<testcase name=\"XUnitTestSuite\\.TestCombine\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<failure message=\"FAIL\" type=\"go.failure\"></failure>\n" +
		" +</testcase>\n" +

		" +<testcase name=\"XUnitTestSuite\\.TestCombine\" classname=\"XUnitTestSuite\" .*file=\"[^\"]*reporter_test.go\".*>\n" +
		" +<skipped>true</skipped>\n" +
		" +</testcase>\n" +

		" +</testsuite>\n" +
		"</testsuites>"

	c.Assert(string(report), Matches, match)
}
