package check_test

import (
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	check "gopkg.in/check.v1"
)

var (
	helperRunFlag = flag.String("helper.run", "", "Run helper suite")
	helperPanicFlag = flag.String("helper.panic", "", "")
)

func TestHelperSuite(t *testing.T) {
	if helperRunFlag == nil || *helperRunFlag == "" {
		t.SkipNow()
	}
	switch *helperRunFlag {
	case "FailHelper":
		check.Run(t, &FailHelper{}, nil)
	case "SuccessHelper":
		check.Run(t, &SuccessHelper{}, nil)
	case "FixtureHelper":
		suite := &FixtureHelper{}
		if helperPanicFlag != nil {
			suite.panicOn = *helperPanicFlag
		}
		check.Run(t, suite, nil)
	default:
		t.Skip()
	}
}

type helperResult []string

var (
	testRunLine = regexp.MustCompile(`^=== (?:RUN|CONT)\s+([0-9A-Za-z/]+)$`)
	testStatusLine = regexp.MustCompile(`^\s*--- ([A-Z]+): ([0-9A-Za-z/]+) \(\d+\.\d+s\)$`)
)

func (result helperResult) Status(test string) string {
	for _, line := range result {
		match := testStatusLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		if match[2] == "TestHelperSuite/" + test {
			return match[1]
		}
	}
	return ""
}

func (result helperResult) Logs(test string) string {
	var lines []string
	var inTest bool
	for _, line := range result {
		if inTest {
			// Log messages are all indented
			if strings.HasPrefix(line, " ") {
				lines = append(lines, line)
				continue
			}
			inTest = false
		}
		match := testRunLine.FindStringSubmatch(line)
		if match != nil && match[1] == "TestHelperSuite/" + test {
			inTest = true
		}
	}
	return strings.Join(lines, "\n")
}

func runHelperSuite(name string, args ...string) (code int, output helperResult) {
	args = append([]string{"-test.v", "-test.run", "TestHelperSuite", "-helper.run", name}, args...)
	cmd := exec.Command(os.Args[0], args...)
	data, err := cmd.Output()
	output = strings.Split(string(data), "\n")
	if execErr, ok := err.(*exec.ExitError); ok {
		code = execErr.ExitCode()
		err = nil
	}
	if err != nil {
		panic(err)
	}
	return code, output
}
