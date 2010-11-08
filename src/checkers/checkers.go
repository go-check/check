package checkers


// Checkers used with the c.Assert() and c.Check() helpers must have
// this interface.
type Checker interface {
    Name() string
    ObtainedLabel() string
    ExpectedLabel() string
    NeedsExpectedValue() bool
    Check(obtained, expected interface{}) bool
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
func (checker *CheckerType) ObtainedLabel() string {
    return "Obtained"
}

// Label to be used for the obtained value when reporting a failure
// in the expectations established.
func (checker *CheckerType) ExpectedLabel() string {
    return "Expected"
}

// Method must return true if the obtained value succeeds the
// expectations established by the given matcher.
func (checker *CheckerType) Check(obtained, expected interface{}) bool {
    return false
}



var Equals Checker = &equalsChecker{}

type equalsChecker struct {
    CheckerType
}

func (checker *equalsChecker) Name() string {
    return "Equals"
}

func (checker *equalsChecker) Check(obtained, expected interface{}) bool {
    return obtained == expected
}
