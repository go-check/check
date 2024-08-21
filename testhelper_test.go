package check_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/check.v1"
)

var (
	helperRunFlag  = flag.String("helper.run", "", "Run helper suite")
	helperFailFlag = flag.String("helper.fail", "", "")
	helperSkipFlag = flag.Int("helper.skip", -1, "")
)

func TestHelperSuite(t *testing.T) {
	if helperRunFlag == nil || *helperRunFlag == "" {
		t.SkipNow()
	}
	switch *helperRunFlag {
	case "FailHelper":
		check.Run(t, &FailHelper{})
	case "SuccessHelper":
		check.Run(t, &SuccessHelper{})
	case "FixtureHelper":
		suite := &FixtureHelper{}
		if helperFailFlag != nil {
			suite.failOn = *helperFailFlag
		}
		if *helperSkipFlag >= 0 {
			suite.skip = true
			suite.skipOnN = *helperSkipFlag
		}
		check.Run(t, suite)
	case "integrationTestHelper":
		check.Run(t, &integrationTestHelper{})
	case "WrongTestArgHelper":
		check.Run(t, &WrongTestArgHelper{})
	case "WrongSetUpTestArgHelper":
		check.Run(t, &WrongSetUpTestArgHelper{})
	case "WrongSetUpSuiteArgHelper":
		check.Run(t, &WrongSetUpSuiteArgHelper{})
	case "WrongTestArgCountHelper":
		check.Run(t, &WrongTestArgCountHelper{})
	case "WrongSetUpTestArgCountHelper":
		check.Run(t, &WrongSetUpTestArgCountHelper{})
	case "WrongSetUpSuiteArgCountHelper":
		check.Run(t, &WrongSetUpSuiteArgCountHelper{})
	case "NoTestsHelper":
		check.Run(t, &NoTestsHelper{})
	case "FixtureCheckHelper":
		suite := &FixtureCheckHelper{}
		if helperFailFlag != nil {
			suite.fail = *helperFailFlag
		}
		check.Run(t, suite)
	case "FixtureLogHelper":
		check.Run(t, &FixtureLogHelper{})
	default:
		t.Skip()
	}
}

// testEvent represents a json formatted event
type testEvent struct {
	Action  string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

type helperResult []testEvent

func (result helperResult) Status(test string) string {
	for _, event := range result {
		if event.Test != "TestHelperSuite/"+test {
			continue
		}
		switch event.Action {
		case "pass", "fail", "skip":
			return event.Action
		}
	}
	return ""
}

var isStatusLine = regexp.MustCompile(`^\s*(?:===|---) `).MatchString

func (result helperResult) logs(match func(testEvent) bool) string {
	var lines []string
	for _, event := range result {
		if !match(event) {
			continue
		}
		if event.Action == "output" && !isStatusLine(event.Output) {
			lines = append(lines, event.Output)
		}
	}
	return strings.Join(lines, "")
}

func (result helperResult) Logs(test string) string {
	return result.logs(func(event testEvent) bool {
		return event.Test == "TestHelperSuite/"+test
	})
}

func (result helperResult) AllLogs() string {
	return result.logs(func(event testEvent) bool {
		return strings.HasPrefix(event.Test, "TestHelperSuite/")
	})
}

func runHelperSuite(c *check.C, name string, args ...string) (code int, output helperResult) {
	c.Helper()
	args = append([]string{"tool", "test2json", os.Args[0], "-test.v", "-test.run", "TestHelperSuite", "-helper.run", name}, args...)
	cmd := exec.Command("go", args...)
	data, err := cmd.Output()
	if execErr, ok := err.(*exec.ExitError); ok {
		code = execErr.ExitCode()
		err = nil
	}
	if err != nil {
		c.Fatal(err)
	}
	lines := bytes.Split(data, []byte("\n"))
	output = make(helperResult, 0, len(lines))
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		var event testEvent
		if err := json.Unmarshal(l, &event); err != nil {
			c.Fatal(err)
		}
		output = append(output, event)
	}
	return code, output
}
