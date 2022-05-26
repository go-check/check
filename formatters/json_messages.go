/*
	Convert go test output to TeamCity format
	Support Run, Skip, Pass, Fail

	For more details see:
		- https://github.com/2tvenom/go-test-teamcity
		- https://confluence.jetbrains.com/display/TCD9/Build+Script+Interaction+with+TeamCity
*/

package formatters

import (
	"encoding/json"
	"strings"
	"time"
)

// -----------------------------------------------------------------------
// Output writer manages atomic output writing according to settings.

type JsonTestEventAction string

// https://github.com/golang/go/blob/master/src/cmd/test2json/main.go
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
func JsonOutput(status string, testName, stdOut string, startTime time.Time, testDuration time.Duration,
	funcPath string) string {
	out := JsonTestEvent{
		Time:    startTime.Format(time.RFC3339),
		Package: funcPath,
		Action:  JsonTestEventActionOutput,
		Test:    testName,
		Elapsed: float64(testDuration) / float64(time.Second),
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
