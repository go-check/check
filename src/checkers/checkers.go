package matchers


// Checkers which are returned by the functions which are used with the
// c.Assert() and c.Check() helpers must match this interface.
type Checker interface {
    FuncName() string
    NeedsExpectedValue() bool
    ObtainedLabel() string
    ExpectedLabel() string
    Check() bool
}

type CheckerFunc func(obtained, expected interface{}) Checker


type CheckerType struct {
    Obtained, Expected interface{}
}


// Trick to ensure it matchers the desired interface.
var _ Checker = (*CheckerType)(nil)


// The function name used to build the matcher.
func (checker *CheckerType) FuncName() string {
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
func (checker *CheckerType) Check() bool {
    return false
}



func Equals(obtained, expected interface{}) Checker {
    return &equalsChecker{CheckerType{obtained, expected}}
}

type equalsChecker struct {
    CheckerType
}

func (checker *equalsChecker) FuncName() string {
    return "Equals"
}

func (checker *equalsChecker) Check() bool {
    return checker.Obtained == checker.Expected
}
