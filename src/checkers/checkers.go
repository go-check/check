package checkers

import (
    "reflect"
    "regexp"
)


// Checkers used with the c.Assert() and c.Check() helpers must have
// this interface.
type Checker interface {
    Name() string
    ObtainedLabel() (varName, varLabel string)
    ExpectedLabel() (varName, varLabel string)
    NeedsExpectedValue() bool
    Check(obtained, expected interface{}) (result bool, error string)
}


type CheckerType struct{}


// Trick to ensure it matchers the desired interface.
var _ Checker = (*CheckerType)(nil)


// The function name used to build the matcher.
func (checker *CheckerType) Name() string {
    return "Checker"
}

// Method must return true if the given matcher needs to be informed
// of an expected value in addition to the actual value obtained to
// verify its expectations.
func (checker *CheckerType) NeedsExpectedValue() bool {
    return true
}

// Label to be used for the obtained value when reporting a failure
// in the expectations established.
func (checker *CheckerType) ObtainedLabel() (varName, varLabel string) {
    return "obtained", "Obtained"
}

// Label to be used for the obtained value when reporting a failure
// in the expectations established.
func (checker *CheckerType) ExpectedLabel() (varName, varLabel string) {
    return "expected", "Expected"
}

// Method must return true if the obtained value succeeds the
// expectations established by the given matcher.  If an error is
// returned, it means the provided parameters are somehow invalid.
func (checker *CheckerType) Check(obtained, expected interface{}) (
        result bool, error string) {
    return false, ""
}



// -----------------------------------------------------------------------
// Equals checker.

var Equals Checker = &equalsChecker{}

type equalsChecker struct {
    CheckerType
}

func (checker *equalsChecker) Name() string {
    return "Equals"
}

func (checker *equalsChecker) Check(obtained, expected interface{}) (
        result bool, error string) {
    // This will use a fast path to check for equality of normal types,
    // and then fallback to reflect.DeepEqual if things go wrong.
    defer func() {
        if recover() != nil {
            result = reflect.DeepEqual(obtained, expected)
        }
    }()
    return obtained == expected, ""
}


// -----------------------------------------------------------------------
// Matches checker.

var Matches Checker = &matchesChecker{}

type matchesChecker struct {
    CheckerType
}

func (checker *matchesChecker) Name() string {
    return "Matches"
}

func (checker *matchesChecker) ObtainedLabel() (varName, varLabel string) {
    return "value", "Value"
}

func (checker *matchesChecker) ExpectedLabel() (varName, varLabel string) {
    return "regex", "Expected to match regex"
}

type hasString interface {
    String() string
}

func (checker *matchesChecker) Check(value, re interface{}) (bool, string) {
    reStr, ok := re.(string)
    if !ok {
        return false, "Regex must be a string"
    }
    valueStr, valueIsStr := value.(string)
    if !valueIsStr {
        if valueWithStr, valueHasStr := value.(hasString); valueHasStr {
            valueStr, valueIsStr = valueWithStr.String(), true
        }
    }
    if valueIsStr {
        matches, err := regexp.MatchString("^" + reStr + "$", valueStr)
        if err != nil {
            return false, "Can't compile regex: " + err.String()
        }
        return matches, ""
    }
    return false, "Obtained value is not a string and has no .String()"
}
