/*
	Convert go test output to TeamCity format
	Support Run, Skip, Pass, Fail

	For more details see:
		- https://github.com/2tvenom/go-test-teamcity
		- https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity
*/

package check

import (
	"encoding/json"
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

func teamcityOutput(status string, testName, stdOut string, startTime time.Time, testDuration time.Duration, details ...string) string {
	now := timeFormat(time.Now())
	testName = escape(*formatMessageNamePrefixFlag + testName)

	if status == "START" {
		return fmt.Sprintf("##teamcity[testStarted timestamp='%s' name='%s' captureStandardOutput='true']",
			startTime.Format(TeamcityTimestampFormat), testName)
	}

	out := ""
	stdOut = strings.TrimSpace(stdOut)
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
		now, testName, testDuration/time.Millisecond)

	return out
}

type JsonTestEventAction string

const (
	JsonTestEventActionRun    = "run"    // the test has started running
	JsonTestEventActionPause  = "pause"  //  the test has been paused
	JsonTestEventActionCount  = "cont"   //  the test has continued running
	JsonTestEventActionPass   = "pass"   //  the test passed
	JsonTestEventActionBench  = "bench"  //  the benchmark printed log output but did not fail
	JsonTestEventActionFail   = "fail"   //  the test or benchmark failed
	JsonTestEventActionOutput = "output" //  the test printed output
	JsonTestEventActionSkip   = "skip"   //  the test was skipped or the package contained no tests
)

type JsonTestEvent struct {
	Time    string // encodes as an RFC3339-format string
	Action  JsonTestEventAction
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

// {"Time":"2022-05-26T13:59:39.562101-04:00","Action":"output","Package":"","Test":"TestService","Output":"OK: 1 passed\n"}
func jsonOutput(status string, testName, stdOut string, startTime time.Time, testDuration time.Duration,
	funcPath, funcName, prefix, label, suffix string) string {
	out := JsonTestEvent{
		Time:    startTime.Format(time.RFC3339),
		Package: fmt.Sprintf("%s%s: %s: %s%s", prefix, label, funcPath, funcName, suffix),
		Action:  JsonTestEventActionOutput,
		Test:    testName,
		Elapsed: 10000 * float64(testDuration) / float64(time.Millisecond),
		Output:  strings.TrimSpace(stdOut),
	}

	switch status {
	case "START":
		out.Action = JsonTestEventActionRun
	case "SKIP":
		out.Action = JsonTestEventActionSkip
	case "MISS":
		out.Action = JsonTestEventActionCount
	case "FAIL":
		out.Action = JsonTestEventActionFail
	case "PASS", "FAIL EXPECTED":
		out.Action = JsonTestEventActionPass
	default: // "PANIC"
		out.Output = JsonTestEventActionFail
	}

	b, _ := json.Marshal(out)
	return string(b)
}
