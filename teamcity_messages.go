/*
	Convert go test output to TeamCity format
	Support Run, Skip, Pass, Fail

	For more details see:
		- https://github.com/2tvenom/go-test-teamcity
		- https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity
*/

package check

import (
	"fmt"
	"strings"
	"time"
)

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

const (
	TeamcityTimestampFormat = "2006-01-02T15:04:05.000"
)

func timeFormat(t time.Time) string {
	return t.Format(TeamcityTimestampFormat)
}

func escapeLines(lines []string) string {
	return escape(strings.Join(lines, "\n"))
}

func escape(s string) string {
	s = strings.Replace(s, "|", "||", -1)
	s = strings.Replace(s, "\n", "|n", -1)
	s = strings.Replace(s, "\r", "|n", -1)
	s = strings.Replace(s, "'", "|'", -1)
	s = strings.Replace(s, "]", "|]", -1)
	s = strings.Replace(s, "[", "|[", -1)
	return s
}

func teamcityOutput(status string, test *C, details ...string) string {
	now := timeFormat(time.Now())
	testName := escape(*formatMessageNamePrefixFlag + test.testName)

	if status == "START" {
		return fmt.Sprintf("##teamcity[testStarted timestamp='%s' name='%s' captureStandardOutput='true']", test.startTime.Format(TeamcityTimestampFormat), testName)
	}

	out := ""
	stdOut := strings.TrimSpace(test.GetTestLog())
	if stdOut != "" {
		out = fmt.Sprintf("##teamcity[testStdOut timestamp='%s' name='%s' out='%s']\n", now, testName, escape(stdOut))
	}

	switch status {
	case "SKIP":
		out += fmt.Sprintf("##teamcity[testIgnored timestamp='%s' name='%s' message='%s']\n", now, testName, "Test is skipped")
	case "MISS":
		out += fmt.Sprintf("##teamcity[testIgnored timestamp='%s' name='%s' message='%s']\n", now, testName, "Test is missed")
	case "FAIL":
		out += fmt.Sprintf("##teamcity[testFailed timestamp='%s' name='%s' details='%s']\n",
			now, testName, escapeLines(details))
	case "PASS", "FAIL EXPECTED":
		// ignore success cases
	default: // "PANIC"
		out += fmt.Sprintf("##teamcity[testFailed timestamp='%s' name='%s' message='Test ended in panic.' details='%s']\n",
			now, testName, escapeLines(details))
	}

	out += fmt.Sprintf("##teamcity[testFinished timestamp='%s' name='%s' duration='%d']",
		now, testName, test.duration/time.Millisecond)

	return out
}
