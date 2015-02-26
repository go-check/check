package check

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
)

// -----------------------------------------------------------------------
// CommentInterface and Commentf helper, to attach extra information to checks.

type comment struct {
	format string
	args   []interface{}
}

// Commentf returns an infomational value to use with Assert or Check calls.
// If the checker test fails, the provided arguments will be passed to
// fmt.Sprintf, and will be presented next to the logged failure.
//
// For example:
//
//     c.Assert(v, Equals, 42, Commentf("Iteration #%d failed.", i))
//
// Note that if the comment is constant, a better option is to
// simply use a normal comment right above or next to the line, as
// it will also get printed with any errors:
//
//     c.Assert(l, Equals, 8192) // Ensure buffer size is correct (bug #123)
//
func Commentf(format string, args ...interface{}) CommentInterface {
	return &comment{format, args}
}

// CommentInterface must be implemented by types that attach extra
// information to failed checks. See the Commentf function for details.
type CommentInterface interface {
	CheckCommentString() string
}

func (c *comment) CheckCommentString() string {
	return fmt.Sprintf(c.format, c.args...)
}

// -----------------------------------------------------------------------
// The Checker interface.

// The Checker interface must be provided by checkers used with
// the Assert and Check verification methods.
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
// Not checker logic inverter.

// The Not checker inverts the logic of the provided checker.  The
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

// The Equals checker verifies that the obtained value is equal to
// the expected value, according to usual Go semantics for ==.
//
// For example:
//
//     c.Assert(value, Equals, 42)
//
var Equals Checker = &equalsChecker{
	&CheckerInfo{Name: "Equals", Params: []string{"obtained", "expected"}},
}

func (checker *equalsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		}
	}()
	return params[0] == params[1], ""
}

// -----------------------------------------------------------------------
// DeepEquals checker.

type deepEqualsChecker struct {
	*CheckerInfo
}

// The DeepEquals checker verifies that the obtained value is deep-equal to
// the expected value.  The check will work correctly even when facing
// slices, interfaces, and values of different types (which always fail
// the test).
//
// For example:
//
//     c.Assert(value, DeepEquals, 42)
//     c.Assert(array, DeepEquals, []string{"hi", "there"})
//
var DeepEquals Checker = &deepEqualsChecker{
	&CheckerInfo{Name: "DeepEquals", Params: []string{"obtained", "expected"}},
}

func (checker *deepEqualsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	return reflect.DeepEqual(params[0], params[1]), ""
}

// -----------------------------------------------------------------------
// HasLen checker.

type hasLenChecker struct {
	*CheckerInfo
}

// The HasLen checker verifies that the obtained value has the
// provided length. In many cases this is superior to using Equals
// in conjuction with the len function because in case the check
// fails the value itself will be printed, instead of its length,
// providing more details for figuring the problem.
//
// For example:
//
//     c.Assert(list, HasLen, 5)
//
var HasLen Checker = &hasLenChecker{
	&CheckerInfo{Name: "HasLen", Params: []string{"obtained", "n"}},
}

func (checker *hasLenChecker) Check(params []interface{}, names []string) (result bool, error string) {
	n, ok := params[1].(int)
	if !ok {
		return false, "n must be an int"
	}
	value := reflect.ValueOf(params[0])
	switch value.Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Chan, reflect.String:
	default:
		return false, "obtained value type has no length"
	}
	if value.Len() == n {
		return true, ""
	}
	return false, fmt.Sprintf("obtained length = %d", value.Len())
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
	params[0] = err.Error()
	names[0] = "error"
	return matches(params[0], params[1])
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
		if valueWithStr, valueHasStr := value.(fmt.Stringer); valueHasStr {
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
}

type doesntPanicChecker struct {
	*CheckerInfo
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
	&CheckerInfo{
		Name:   "Panics",
		Params: []string{"function", "expected"}},
}

var DoesntPanic Checker = &doesntPanicChecker{
	&CheckerInfo{Name: "DoesntPanic", Params: []string{"function"}},
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
		params[0] = recover()
		names[0] = "panic"
		result = reflect.DeepEqual(params[0], params[1])
	}()
	f.Call(nil)
	return false, "Function has not panicked"
}

// The DoesntPanic checker verifies that calling the provided zero-argument
// function will NOT cause a panic.
//
// The first param must be a function so that the execution of the code
// to be tested can be delayed, and any unexpected panic caught.

// For example:
//
//     c.Assert(func() { f(1, 2) }, DoesntPanic)
//
//
func (checker *doesntPanicChecker) Check(params []interface{}, names []string) (result bool, err string) {
	f := reflect.ValueOf(params[0])
	if f.Kind() != reflect.Func || f.Type().NumIn() != 0 {
		return false, "Function must take zero arguments"
	}
	defer func() {
		if e := recover(); e != nil {
			result = false
			// TODO: Figure out how to set the error string
			err = fmt.Sprintf("%v", e)
		}
	}()
	f.Call(nil)
	return true, ""
}

type panicMatchesChecker struct {
	*CheckerInfo
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
var PanicMatches Checker = &panicMatchesChecker{
	&CheckerInfo{Name: "PanicMatches", Params: []string{"function", "expected"}},
}

func (checker *panicMatchesChecker) Check(params []interface{}, names []string) (result bool, errmsg string) {
	f := reflect.ValueOf(params[0])
	if f.Kind() != reflect.Func || f.Type().NumIn() != 0 {
		return false, "Function must take zero arguments"
	}
	defer func() {
		// If the function has not panicked, then don't do the check.
		if errmsg != "" {
			return
		}
		obtained := recover()
		names[0] = "panic"
		if e, ok := obtained.(error); ok {
			params[0] = e.Error()
		} else if _, ok := obtained.(string); ok {
			params[0] = obtained
		} else {
			errmsg = "Panic value is not a string or an error"
			return
		}
		result, errmsg = matches(params[0], params[1])
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

// -----------------------------------------------------------------------
// IsTrue / IsFalse checker.

type isBoolValueChecker struct {
	*CheckerInfo
	expected bool
}

func (checker *isBoolValueChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained, ok := params[0].(bool)
	if !ok {
		return false, "Argument to " + checker.Name + " must be bool"
	}

	return obtained == checker.expected, ""
}

// The IsTrue checker verifies that the obtained value is true.
//
// For example:
//
//     c.Assert(value, IsTrue)
//
var IsTrue Checker = &isBoolValueChecker{
	&CheckerInfo{Name: "IsTrue", Params: []string{"obtained"}},
	true,
}

// The IsFalse checker verifies that the obtained value is false.
//
// For example:
//
//     c.Assert(value, IsFalse)
//
var IsFalse Checker = &isBoolValueChecker{
	&CheckerInfo{Name: "IsFalse", Params: []string{"obtained"}},
	false,
}

// -----------------------------------------------------------------------
// SliceIncludes checker.

type sliceIncludesChecker struct {
	*CheckerInfo
}

// The SliceIncludes checker verifies that the provided slice includes
// the provided object.
//
// For example:
// c.Assert(aSlice, SliceIncludes, aThing)
//
var SliceIncludes Checker = &sliceIncludesChecker{
	&CheckerInfo{
		Name:   "SliceIncludes",
		Params: []string{"aSlice", "aThing"}},
}

func (checker *sliceIncludesChecker) Check(params []interface{}, names []string) (result bool, errStr string) {
	//params[0] == aSlice
	//params[1] == aThing (that we hope is in the slice)
	s := reflect.ValueOf(params[0]) //aSlice
	if s.Kind() != reflect.Slice {
		return false, fmt.Sprintf("SliceIncludes given a non-slice type: %v", params[0])
	}

	for i := 0; i < s.Len(); i++ {
		x := s.Index(i).Interface()
		if params[1] == x { //roughly:  if aThing == aSlice[i]
			return true, ""
		}
	}
	return false, ""

}

// The WithinDelta checker verifies that the obtained value is
// withen a delta of the expected value. All numbers must be floats64s
//
// For example:
// c.Assert(gear.GearInches(), WithinDelta, 0.01, 137.1)//
//
// Based on Nathan Youngman's <http://nathany.com/> work
// found here:
// <https://github.com/nathany/go-poodr/blob/master/chapter9/gear1/gear1_check_test.go>
// Distributed under the simplified BSD license:
// <https://github.com/nathany/go-poodr/blob/master/LICENSE>
type withinDeltaChecker struct {
	*CheckerInfo
}

var WithinDelta Checker = &withinDeltaChecker{
	&CheckerInfo{
		Name:   "WithinDelta",
		Params: []string{"obtained", "delta", "expected"}},
}

func (c *withinDeltaChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained, ok := params[0].(float64)
	if !ok {
		return false, "obtained must be a float64"
	}
	delta, ok := params[1].(float64)
	if !ok {
		return false, "delta must be a float64"
	}
	expected, ok := params[2].(float64)
	if !ok {
		return false, "expected must be a float64"
	}
	return math.Abs(obtained-expected) <= delta, ""
}

// BetweenFloats Checker
// checks if a float is between a low and high value
// similar to WithinDelta but the range on each side isn't
// necessarily balanced.
//
// c.Assert(child.Age(), BetweenFloats, 3.5, 5.0)
type betweenFloatsChecker struct {
	*CheckerInfo
}

var BetweenFloats Checker = &betweenFloatsChecker{
	&CheckerInfo{
		Name:   "BetweenFloats",
		Params: []string{"obtained", "low", "high"}},
}

func (c *betweenFloatsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained, ok := params[0].(float64)
	if !ok {
		return false, "obtained must be a float64"
	}
	low, ok := params[1].(float64)
	if !ok {
		return false, "low must be a float64"
	}
	high, ok := params[2].(float64)
	if !ok {
		return false, "high must be a float64"
	}
	return (obtained >= low && obtained <= high), ""
}
