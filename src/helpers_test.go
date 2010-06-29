// These tests verify the inner workings of the helper methods associated
// with gocheck.T.

package gocheck_test

import (
    "gocheck"
    "testing"
    "fmt"
)


func TestHelpers(t *testing.T) {
    gocheck.RunTestingT(&HelpersS{}, t)
}


// -----------------------------------------------------------------------
// Helpers test suite.

type HelpersS struct{}

func (s *HelpersS) TestCheckEqualSucceeding(t *gocheck.T) {
    result := t.CheckEqual(10, 10)
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 10)",
                    result: true,
               })
}

func (s *HelpersS) TestCheckEqualFailing(t *gocheck.T) {
    log := fmt.Sprintf(
        "\n... %d:CheckEqual(A, B): A != B\n" +
        "... A: 10\n" +
        "... B: 20\n\n",
        getMyLine()+1)
    result := t.CheckEqual(10, 20)
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 20)",
                    result: false,
                    failed: true,
                    log: log,
               })
}

func (s *HelpersS) TestCheckEqualWithMessage(t *gocheck.T) {
    log := fmt.Sprintf(
        "\n... %d:CheckEqual(A, B): A != B\n" +
        "... A: 10\n" +
        "... B: 20\n" +
        "... That's clearly.. WRONG!\n\n",
        getMyLine()+1)
    result := t.CheckEqual(10, 20, "That's clearly.. ", "WRONG!")
    checkState(t, result,
               &expectedState{
                    name: "CheckEqual(10, 20)",
                    result: false,
                    failed: true,
                    log: log,
               })
}

func (s *HelpersS) TestCheckNotEqualSucceeding(t *gocheck.T) {
    result := t.CheckNotEqual(10, 20)
    checkState(t, result,
               &expectedState{
                    name: "CheckNotEqual(10, 20)",
                    result: true,
               })
}

func (s *HelpersS) TestCheckNotEqualFailing(t *gocheck.T) {
    log := fmt.Sprintf(
        "\n... %d:CheckNotEqual(A, B): A == B\n" +
        "... A: 10\n" +
        "... B: 10\n\n",
        getMyLine()+1)
    result := t.CheckNotEqual(10, 10)
    checkState(t, result,
               &expectedState{
                    name: "CheckNotEqual(10, 10)",
                    result: false,
                    failed: true,
                    log: log,
               })
}

func (s *HelpersS) TestCheckNotEqualWithMessage(t *gocheck.T) {
    log := fmt.Sprintf(
        "\n... %d:CheckNotEqual(A, B): A == B\n" +
        "... A: 10\n" +
        "... B: 10\n" +
        "... That's clearly.. WRONG!\n\n",
        getMyLine()+1)
    result := t.CheckNotEqual(10, 10, "That's clearly.. ", "WRONG!")
    checkState(t, result,
               &expectedState{
                    name: "CheckNotEqual(10, 10)",
                    result: false,
                    failed: true,
                    log: log,
               })
}
