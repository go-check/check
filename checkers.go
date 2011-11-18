package gocheck

import (
	"reflect"
	"regexp"
	"fmt"
)

// -----------------------------------------------------------------------
// BugInfo and Bug() helper, to attach extra information to checks.

type bugInfo struct {
	format string
	args   []interface{}
}

// Bug enables attaching some information to Assert() or Check() calls.
// If the checker test fails, the provided arguments will be passed to
// fmt.Sprintf(), and will be presented next to the logged failure.
//
// For example:
//
//     c.Assert(l, Equals, 8192, Bug("Buffer size is incorrect, bug #123"))
//     c.Assert(v, Equals, 42, Bug("Iteration #%d", i))
//
func Bug(format string, args ...interface{}) BugInfo {
	return &bugInfo{format, args}
}

// BugInfo is the interface which must be supported for attaching extra
// information to checks.  See the Bug() function for details.
type BugInfo interface {
	GetBugInfo() string
}

func (bug *bugInfo) GetBugInfo() string {
	return fmt.Sprintf(bug.format, bug.args...)
}

// -----------------------------------------------------------------------
// The Checker interface.

// The Checker interface must be provided by checkers used with
// the c.Assert() and c.Check() verification methods.
type Checker interface {
	Info() *CheckerInfo
	Check(params []interface{}, names []string) (result bool, error string)
}

// See the Checker interface.
type CheckerInfo struct {
	Name   string
	Params []string
}

func (info *CheckerInfo) Info() *CheckerInfo {
	return info
}

// -----------------------------------------------------------------------
// Not() checker logic inverter.

// The Not() checker inverts the logic of the provided checker.  The
// resulting checker will succeed where the original one failed, and
// vice-versa.
//
// For example:
//
//     c.Assert(a, Not(Equals), b)
//
func Not(checker Checker) Checker {
	return &notChecker{checker}
}

type notChecker struct {
	sub Checker
}

func (checker *notChecker) Info() *CheckerInfo {
	info := *checker.sub.Info()
	info.Name = "Not(" + info.Name + ")"
	return &info
}

func (checker *notChecker) Check(params []interface{}, names []string) (result bool, error string) {
	result, error = checker.sub.Check(params, names)
	result = !result
	return
}

// -----------------------------------------------------------------------
// IsNil checker.

type isNilChecker struct {
	*CheckerInfo
}

// The IsNil checker tests whether the obtained value is nil.
//
// For example:
//
//    c.Assert(err, IsNil)
//
var IsNil Checker = &isNilChecker{
	&CheckerInfo{Name: "IsNil", Params: []string{"value"}},
}

func (checker *isNilChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return isNil(params[0]), ""
}

func isNil(obtained interface{}) (result bool) {
	if obtained == nil {
		result = true
	} else {
		switch v := reflect.ValueOf(obtained); v.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			return v.IsNil()
		}
	}
	return
}

// -----------------------------------------------------------------------
// NotNil checker. Alias for Not(IsNil), since it's so common.

type notNilChecker struct {
	*CheckerInfo
}

// The NotNil checker verifies that the obtained value is not nil.
//
// For example:
//
//     c.Assert(iface, NotNil)
//
// This is an alias for Not(IsNil), made available since it's a
// fairly common check.
//
var NotNil Checker = &notNilChecker{
	&CheckerInfo{Name: "NotNil", Params: []string{"value"}},
}

func (checker *notNilChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return !isNil(params[0]), ""
}

// -----------------------------------------------------------------------
// Equals checker.

type equalsChecker struct {
	*CheckerInfo
}

// The Equals checker verifies that the obtained value is deep-equal to
// the expected value.  The check will work correctly even when facing
// arrays, interfaces, and values of different types (which always fail
// the test).
//
// For example:
//
//     c.Assert(value, Equals, 42)
//     c.Assert(array, Equals, []string{"hi", "there"})
//
var Equals Checker = &equalsChecker{
	&CheckerInfo{Name: "Equals", Params: []string{"obtained", "expected"}},
}

func (checker *equalsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return reflect.DeepEqual(params[0], params[1]), ""
}

// -----------------------------------------------------------------------
// ErrorMatches checker.

type errorMatchesChecker struct {
	*CheckerInfo
}

// The ErrorMatches checker verifies that the error value
// is non nil and matches the regular expression provided.
//
// For example:
//
//     c.Assert(err, ErrorMatches, "perm.*denied")
//
var ErrorMatches Checker = errorMatchesChecker{
	&CheckerInfo{Name: "ErrorMatches", Params: []string{"value", "regex"}},
}

func (checker errorMatchesChecker) Check(params []interface{}, names []string) (result bool, errStr string) {
	if params[0] == nil {
		return false, "Error value is nil"
	}
	err, ok := params[0].(error)
	if !ok {
		return false, "Value is not an error"
	}
	return matches(err.Error(), params[1])
}

// -----------------------------------------------------------------------
// Matches checker.

type matchesChecker struct {
	*CheckerInfo
}

// The Matches checker verifies that the string provided as the obtained
// value (or the string resulting from obtained.String()) matches the
// regular expression provided.
//
// For example:
//
//     c.Assert(err, Matches, "perm.*denied")
//
var Matches Checker = &matchesChecker{
	&CheckerInfo{Name: "Matches", Params: []string{"value", "regex"}},
}

func (checker *matchesChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return matches(params[0], params[1])
}

func matches(value, regex interface{}) (result bool, error string) {
	reStr, ok := regex.(string)
	if !ok {
		return false, "Regex must be a string"
	}
	valueStr, valueIsStr := value.(string)
	if !valueIsStr {
		if valueWithStr, valueHasStr := value.(stringer); valueHasStr {
			valueStr, valueIsStr = valueWithStr.String(), true
		}
	}
	if valueIsStr {
		matches, err := regexp.MatchString("^"+reStr+"$", valueStr)
		if err != nil {
			return false, "Can't compile regex: " + err.Error()
		}
		return matches, ""
	}
	return false, "Obtained value is not a string and has no .String()"
}

// -----------------------------------------------------------------------
// Panics checker.

type panicsChecker struct {
	*CheckerInfo
	check func(obtained, expected interface{}) (bool, string)
}

// The Panics checker verifies that calling the provided zero-argument
// function will cause a panic which is deep-equal to the provided value.
//
// For example:
//
//     c.Assert(func() { f(1, 2) }, Panics, &SomeErrorType{"BOOM"}).
//
//
var Panics Checker = &panicsChecker{
	&CheckerInfo{Name: "Panics", Params: []string{"function", "expected"}},
	func(obtained, expected interface{}) (bool, string) {
		return reflect.DeepEqual(obtained, expected), ""
	},
}

// The PanicMatches checker verifies that calling the provided zero-argument
// function will cause a panic with an error value matching
// the regular expression provided.
//
// For example:
//
//     c.Assert(func() { f(1, 2) }, PanicMatches, `open.*: no such file or directory`).
//
//
var PanicMatches Checker = &panicsChecker{
	&CheckerInfo{Name: "PanicMatches", Params: []string{"function", "expected"}},
	func(obtained, expected interface{}) (bool, string) {
		v, ok := obtained.(error)
		if !ok {
			return false, "Panic value is not an error"
		}
		return matches(v.Error(), expected)
	},
}

func (checker *panicsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	f := reflect.ValueOf(params[0])
	if f.Kind() != reflect.Func || f.Type().NumIn() != 0 {
		return false, "Function must take zero arguments"
	}
	defer func() {
		// If the function has not panicked, then don't do the check.
		if error != "" {
			return
		}
		obtained := recover()
		expected := params[1]
		params[0] = obtained
		names[0] = "panic"

		result, error = checker.check(obtained, expected)
	}()
	f.Call(nil)
	return false, "Function has not panicked"
}

// -----------------------------------------------------------------------
// FitsTypeOf checker.

type fitsTypeChecker struct {
	*CheckerInfo
}

// The FitsTypeOf checker verifies that the obtained value is
// assignable to a variable with the same type as the provided
// sample value.
//
// For example:
//
//     c.Assert(value, FitsTypeOf, int64(0))
//     c.Assert(value, FitsTypeOf, os.Error(nil))
//
var FitsTypeOf Checker = &fitsTypeChecker{
	&CheckerInfo{Name: "FitsTypeOf", Params: []string{"obtained", "sample"}},
}

func (checker *fitsTypeChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained := reflect.ValueOf(params[0])
	sample := reflect.ValueOf(params[1])
	if !obtained.IsValid() {
		return false, ""
	}
	if !sample.IsValid() {
		return false, "Invalid sample value"
	}
	return obtained.Type().AssignableTo(sample.Type()), ""
}

// -----------------------------------------------------------------------
// Implements checker.

type implementsChecker struct {
	*CheckerInfo
}

// The Implements checker verifies that the obtained value
// implements the interface specified via a pointer to an interface
// variable.
//
// For example:
//
//     var e os.Error
//     c.Assert(err, Implements, &e)
//
var Implements Checker = &implementsChecker{
	&CheckerInfo{Name: "Implements", Params: []string{"obtained", "ifaceptr"}},
}

func (checker *implementsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained := reflect.ValueOf(params[0])
	ifaceptr := reflect.ValueOf(params[1])
	if !obtained.IsValid() {
		return false, ""
	}
	if !ifaceptr.IsValid() || ifaceptr.Kind() != reflect.Ptr || ifaceptr.Elem().Kind() != reflect.Interface {
		return false, "ifaceptr should be a pointer to an interface variable"
	}
	return obtained.Type().Implements(ifaceptr.Elem().Type()), ""
}
