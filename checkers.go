package check

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"

	cf "github.com/iostrovok/go-convert"
	"github.com/niemeyer/pretty"
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
	if result {
		// clear error message if the new result is true
		error = ""
	}
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

func diffworthy(a interface{}) bool {
	if a == nil {
		return false
	}

	t := reflect.TypeOf(a)
	switch t.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct, reflect.String, reflect.Ptr:
		return true
	}
	return false
}

// formatUnequal will dump the actual and expected values into a textual
// representation and return an error message containing a diff.
func formatUnequal(obtained interface{}, expected interface{}) string {
	// We do not do diffs for basic types because go-check already
	// shows them very cleanly.
	if !diffworthy(obtained) || !diffworthy(expected) {
		return ""
	}

	// Handle strings, short strings are ignored (go-check formats
	// them very nicely already). We do multi-line strings by
	// generating two string slices and using kr.Diff to compare
	// those (kr.Diff does not do string diffs by itself).
	aStr, aOK := obtained.(string)
	bStr, bOK := expected.(string)
	if aOK && bOK {
		l1 := strings.Split(aStr, "\n")
		l2 := strings.Split(bStr, "\n")
		// the "2" here is a bit arbitrary
		if len(l1) > 2 && len(l2) > 2 {
			diff := pretty.Diff(l1, l2)
			return fmt.Sprintf(`String difference:
%s`, formatMultiLine(strings.Join(diff, "\n"), false))
		}
		// string too short
		return ""
	}

	// generic diff
	diff := pretty.Diff(obtained, expected)
	if len(diff) == 0 {
		// No diff, this happens when e.g. just struct
		// pointers are different but the structs have
		// identical values.
		return ""
	}

	return fmt.Sprintf(`Difference:
%s`, formatMultiLine(strings.Join(diff, "\n"), false))
}

func formatUnsupportedType(params []interface{}) string {
	out := "Comparing incomparable type " +
		reflect.ValueOf(params[0]).Type().String() +
		" and " +
		reflect.ValueOf(params[1]).Type().String()

	return out
}

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

	result = params[0] == params[1]
	if !result {
		error = formatUnequal(params[0], params[1])
	}
	return
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
	result = reflect.DeepEqual(params[0], params[1])
	if !result {
		error = formatUnequal(params[0], params[1])
	}
	return
}

// -----------------------------------------------------------------------
// HasLen checker.

type hasLenChecker struct {
	*CheckerInfo
}

// The HasLen checker verifies that the obtained value has the
// provided length. In many cases this is superior to using Equals
// in conjunction with the len function because in case the check
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
	n := 0
	switch reflect.ValueOf(params[1]).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = cf.Int(params[1])
	default:
		return false, fmt.Sprintf("n must be an int*, not %T", params[1])
	}

	value := reflect.ValueOf(params[0])
	switch value.Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Chan, reflect.String:
	default:
		return false, "obtained value type has no length property"
	}
	return value.Len() == n, ""
}

// -----------------------------------------------------------------------
// HasLenMoreThan checker.

type hasLenMoreThan struct {
	*CheckerInfo
}

// The HasLenMoreThan checker verifies that the obtained value has
// the length more than provided one. In many cases this is superior
// to using Equals in conjunction with the len() function because
// in case the check fails the value itself will be printed, instead of its length,
// providing more details for figuring the problem.
// Also it converts last parameter from any int* to int, that may be useful for use configured/calculated values.
//
// For example:
//
//     c.Assert(list, HasLenMoreThan, 5)
//
var HasLenMoreThan Checker = &hasLenMoreThan{
	&CheckerInfo{Name: "HasLenMoreThan", Params: []string{"obtained", "n"}},
}

func (checker *hasLenMoreThan) Check(params []interface{}, names []string) (result bool, error string) {
	n := 0
	switch reflect.ValueOf(params[1]).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = cf.Int(params[1])
	default:
		return false, fmt.Sprintf("n must be an int*, not %T", params[1])
	}

	value := reflect.ValueOf(params[0])
	switch value.Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Chan, reflect.String:
	default:
		return false, "obtained value type has no length property"
	}
	return value.Len() > n, ""
}

// -----------------------------------------------------------------------
// HasLenLessThan checker.l

type hasLenLessThan struct {
	*CheckerInfo
}

// The hasLenLessThan checker verifies that the obtained value has
// the length less than provided one. In many cases this is superior
// to using Equals in conjunction with the len() function because
// in case the check fails the value itself will be printed, instead of its length,
// providing more details for figuring the problem.
// Also it converts last parameter from any int* to int, that may be useful for use configured/calculated values.
//
// For example:
//
//     c.Assert(list, HasLenLessThan, 5)
//
var HasLenLessThan Checker = &hasLenLessThan{
	&CheckerInfo{Name: "HasLenLessThan", Params: []string{"obtained", "n"}},
}

func (checker *hasLenLessThan) Check(params []interface{}, names []string) (result bool, error string) {
	n := 0
	switch reflect.ValueOf(params[1]).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = cf.Int(params[1])
	default:
		return false, fmt.Sprintf("n must be an int*, not %T", params[1])
	}

	value := reflect.ValueOf(params[0])
	switch value.Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Chan, reflect.String:
	default:
		return false, "obtained value type has no length property"
	}
	return value.Len() < n, ""
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
// MoreThan checker.

var NoMoreThanStringError = "Difference: first (string) parameter <= (string) second, expect more"

type moreThan struct {
	*CheckerInfo
}

// The MoreThan checker tests whether the obtained value is more than value.
//
// For example:
//
//    c.Assert(v1, MoreThan, v2) -> v1 > v2
//    c.Assert(v1, LessThan, 23) -> v1 > 23
//    c.Assert(v1, LessThan, "my string") -> v1 > "my string"
//    c.Assert(v1, LessThan, []byte("my string")) -> string(v1) > "my string"
//
//	  defaults conversion for checking:
//    Int, Int8, Int16, Int32, Int64 => Int64
//	  Uint, Uint8, Uint16, Uint32, Uint64 => Uint64
//	  float32 => float32
//    []byte, string => string
//    float64 => float64
//

var MoreThan Checker = &moreThan{
	&CheckerInfo{Name: "MoreThan", Params: []string{"get", "should be more"}},
}

func (checker *moreThan) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		} else if !result && error == "" {
			error = fmt.Sprintf("Difference: %v <= %v", params[0], params[1])
		}
	}()

	if a := []bool{isStringType(params[0]), isStringType(params[1])}; a[0] || a[1] {
		if a[0] && a[1] {
			if result = cf.String(params[0]) > cf.String(params[1]); !result {
				error = "First (string) parameter equals (string) second, expect more"

				// generic diff
				if diff := pretty.Diff(cf.String(params[0]), cf.String(params[1])); len(diff) > 0 {
					error = NoLessThanStringError
				}
			}
		} else {
			error = formatUnsupportedType(params)
		}
		return
	}

	rt := reflect.ValueOf(params[0]).Kind()
	if rt == reflect.Ptr {
		rt = reflect.ValueOf(params[0]).Type().Kind()
	}

	rt2 := reflect.ValueOf(params[1]).Kind()
	if rt2 == reflect.Ptr {
		rt2 = reflect.ValueOf(params[1]).Type().Kind()
	}

	// unsupported types
	if rt != rt2 {
		if reflect.Float32 == rt || reflect.Float32 == rt2 || reflect.Float64 == rt || reflect.Float64 == rt2 {
			error = formatUnsupportedType(params)
			return
		}
	}

	switch rt {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = cf.Int64(params[0]) > cf.Int64(params[1])
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = cf.Uint64(params[0]) > cf.Uint64(params[1])
	case reflect.Float32:
		result = cf.Float32(params[0]) > cf.Float32(params[1])
	case reflect.Float64:
		result = cf.Float64(params[0]) > cf.Float64(params[1])
	default:
		error = formatUnsupportedType(params)
	}

	return
}

// -----------------------------------------------------------------------
// LessThan checker.

var NoLessThanStringError = "Difference: first (string) parameter <= (string) second, expect less"

type lessThan struct {
	*CheckerInfo
}

// The LessThan checker tests whether the obtained value is less than value.
//
// For example:
//
//    c.Assert(v1, LessThan, v2) -> v1 < v2
//    c.Assert(v1, LessThan, 23) -> v1 < 23
//    c.Assert(v1, LessThan, "my string") -> v1 < "my string"
//    c.Assert(v1, LessThan, []byte("my string")) -> string(v1) < "my string"
//
//	  defaults conversion for checking:
//    Int, Int8, Int16, Int32, Int64 => Int64
//	  Uint, Uint8, Uint16, Uint32, Uint64 => Uint64
//	  float32 => float32
//    []byte, string => string
//    float64 => float64
//

var LessThan Checker = &lessThan{
	&CheckerInfo{Name: "LessThan", Params: []string{"get", "should be less"}},
}

func (checker *lessThan) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		} else if !result && error == "" {
			error = fmt.Sprintf("Difference: %v >= %v", params[0], params[1])
		}
	}()

	if a := []bool{isStringType(params[0]), isStringType(params[1])}; a[0] || a[1] {
		if a[0] && a[1] {
			result = cf.String(params[0]) < cf.String(params[1])
			if !result {
				// generic diff
				if diff := pretty.Diff(cf.String(params[0]), cf.String(params[1])); len(diff) > 0 {
					error = NoLessThanStringError
				} else {
					error = "First (string) parameter equals (string) second, expect less"
				}
			}
		} else {
			error = formatUnsupportedType(params)
		}
		return
	}

	rt := reflect.ValueOf(params[0]).Kind()
	if rt == reflect.Ptr {
		rt = reflect.ValueOf(params[0]).Type().Kind()
	}

	rt2 := reflect.ValueOf(params[1]).Kind()
	if rt2 == reflect.Ptr {
		rt2 = reflect.ValueOf(params[1]).Type().Kind()
	}

	// unsupported types
	if rt != rt2 {
		if reflect.Float32 == rt || reflect.Float32 == rt2 || reflect.Float64 == rt || reflect.Float64 == rt2 {
			error = formatUnsupportedType(params)
			return
		}
	}

	switch rt {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = cf.Int64(params[0]) < cf.Int64(params[1])
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = cf.Uint64(params[0]) < cf.Uint64(params[1])
	case reflect.Float32:
		result = cf.Float32(params[0]) < cf.Float32(params[1])
	case reflect.Float64:
		result = cf.Float64(params[0]) < cf.Float64(params[1])
	default:
		error = formatUnsupportedType(params)
	}

	return
}

// -----------------------------------------------------------------------
// MoreOrEqualThan checker.

var NoMoreOrEqualThanStringError = "Difference: first (string) parameter > (string) second, expect less or equal"

type moreOrEqualThan struct {
	*CheckerInfo
}

// The MoreOrEqualThan checker tests whether the obtained value is less than value.
//
// For example:
//
//    c.Assert(v1, MoreOrEqualThan, v2) -> v1 = v2
//    c.Assert(v1, MoreOrEqualThan, 23) -> v1 >= 23
//    c.Assert(v1, MoreOrEqualThan, "my string") -> v1 >= "my string"
//    c.Assert(v1, MoreOrEqualThan, []byte("my string")) -> string(v1) >= "my string"
//
//	  defaults conversion for checking:
//    Int, Int8, Int16, Int32, Int64 => Int64
//	  Uint, Uint8, Uint16, Uint32, Uint64 => Uint64
//	  float32 => float32
//    []byte, string => string
//    float64 => float64
//

var MoreOrEqualThan Checker = &moreOrEqualThan{
	&CheckerInfo{Name: "MoreOrEqualThan", Params: []string{"get", "should be more or equal than"}},
}

func (checker *moreOrEqualThan) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		} else if !result && error == "" {
			error = fmt.Sprintf("Difference: %v < %v", params[0], params[1])
		}
	}()

	if a := []bool{isStringType(params[0]), isStringType(params[1])}; a[0] || a[1] {
		if a[0] && a[1] {
			if result = cf.String(params[0]) >= cf.String(params[1]); !result {
				error = NoMoreOrEqualThanStringError
			}
		} else {
			error = formatUnsupportedType(params)
		}
		return
	}

	rt := reflect.ValueOf(params[0]).Kind()
	if rt == reflect.Ptr {
		rt = reflect.ValueOf(params[0]).Type().Kind()
	}

	rt2 := reflect.ValueOf(params[1]).Kind()
	if rt2 == reflect.Ptr {
		rt2 = reflect.ValueOf(params[1]).Type().Kind()
	}

	// unsupported types
	if rt != rt2 {
		if reflect.Float32 == rt || reflect.Float32 == rt2 || reflect.Float64 == rt || reflect.Float64 == rt2 {
			error = formatUnsupportedType(params)
			return
		}
	}

	switch rt {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = cf.Int64(params[0]) >= cf.Int64(params[1])
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = cf.Uint64(params[0]) >= cf.Uint64(params[1])
	case reflect.Float32:
		result = cf.Float32(params[0]) >= cf.Float32(params[1])
	case reflect.Float64:
		result = cf.Float64(params[0]) >= cf.Float64(params[1])
	default:
		error = formatUnsupportedType(params)
	}

	return
}

func isStringType(p interface{}) bool {
	strName := reflect.ValueOf(p).Type().String()
	return strName == "string" || strName == "[]uint8"
}

// -----------------------------------------------------------------------
// LessOrEqualThan checker.
var NoLessOrEqualThanStringError = "Difference: first (string) parameter > (string) second, expect less or equal"

type lessOrEqualThan struct {
	*CheckerInfo
}

// The LessOrEqualThan checker tests whether the obtained value is less than value.
//
// For example:
//
//    c.Assert(v1, LessOrEqualThan, 23) -> v1 <= v2
//    c.Assert(v1, LessOrEqualThan, "my string") -> v1 <= v2
//
//	  defaults conversion for checking:
//    Int, Int8, Int16, Int32, Int64 => Int64
//	  Uint, Uint8, Uint16, Uint32, Uint64 => Uint64
//	  float32 => float32
//    []byte, string => string
//    float64 => float64
//
var LessOrEqualThan Checker = &lessOrEqualThan{
	&CheckerInfo{Name: "MoreOrEqualThan", Params: []string{"get", "should be more or equal than"}},
}

func (checker *lessOrEqualThan) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		} else if !result && error == "" {
			error = fmt.Sprintf("Difference: %v > %v", params[0], params[1])
		}
	}()

	if a := []bool{isStringType(params[0]), isStringType(params[1])}; a[0] || a[1] {
		if a[0] && a[1] {
			if result = cf.String(params[0]) <= cf.String(params[1]); !result {
				error = NoLessOrEqualThanStringError
			}
		} else {
			error = formatUnsupportedType(params)
		}
		return
	}

	rt := reflect.ValueOf(params[0]).Kind()
	if rt == reflect.Ptr {
		rt = reflect.ValueOf(params[0]).Type().Kind()
	}

	rt2 := reflect.ValueOf(params[1]).Kind()
	if rt2 == reflect.Ptr {
		rt2 = reflect.ValueOf(params[1]).Type().Kind()
	}

	// unsupported types
	if rt != rt2 {
		if reflect.Float32 == rt || reflect.Float32 == rt2 || reflect.Float64 == rt || reflect.Float64 == rt2 {
			error = formatUnsupportedType(params)
			return
		}
	}

	switch rt {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = cf.Int64(params[0]) <= cf.Int64(params[1])
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = cf.Uint64(params[0]) <= cf.Uint64(params[1])
	case reflect.Float32:
		result = cf.Float32(params[0]) <= cf.Float32(params[1])
	case reflect.Float64:
		result = cf.Float64(params[0]) <= cf.Float64(params[1])
	default:
		error = formatUnsupportedType(params)
	}

	return
}

// -----------------------------------------------------------------------
// EqualsMore checker.
var NoEqualsMoreStringError = "Difference: first (string) parameter does not equals (string) second, expect 'equals'"

type equalsMore struct {
	*CheckerInfo
}

// The EqualsMore checker verifies that the obtained value is equal to
// the expected value, as Equals checker, but converts parameters if it is possible.
//
// For example:
//
//     c.Assert(value, Equals, 42)
//     c.Assert(42, Equals, int64(42)) // true
//     c.Assert(int32(42), Equals, int64(42)) // true
//
//	  defaults conversion for checking:
//    Int, Int8, Int16, Int32, Int64 => Int64
//	  Uint, Uint8, Uint16, Uint32, Uint64 => Uint64
//	  float32 => float32
//    []byte, string => string
//    float64 => float64
//

var EqualsMore Checker = &equalsMore{
	&CheckerInfo{Name: "EqualsMore", Params: []string{"get", "should be equal"}},
}

func (checker *equalsMore) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		} else if !result && error == "" {
			error = fmt.Sprintf("Difference: %v != %v", params[0], params[1])
		}
	}()

	if a := []bool{isStringType(params[0]), isStringType(params[1])}; a[0] || a[1] {
		if a[0] && a[1] {
			if result = cf.String(params[0]) == cf.String(params[1]); !result {
				error = NoEqualsMoreStringError
			}
		} else {
			error = formatUnsupportedType(params)
		}
		return
	}

	rt := reflect.ValueOf(params[0]).Kind()
	if rt == reflect.Ptr {
		rt = reflect.ValueOf(params[0]).Type().Kind()
	}

	rt2 := reflect.ValueOf(params[1]).Kind()
	if rt2 == reflect.Ptr {
		rt2 = reflect.ValueOf(params[1]).Type().Kind()
	}

	// unsupported types
	if rt != rt2 {
		if reflect.Float32 == rt || reflect.Float32 == rt2 || reflect.Float64 == rt || reflect.Float64 == rt2 {
			error = formatUnsupportedType(params)
			return
		}
	}

	switch rt {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = cf.Int64(params[0]) == cf.Int64(params[1])
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = cf.Uint64(params[0]) == cf.Uint64(params[1])
	case reflect.Float32:
		result = cf.Float32(params[0]) == cf.Float32(params[1])
	case reflect.Float64:
		result = cf.Float64(params[0]) == cf.Float64(params[1])
	default:
		error = formatUnsupportedType(params)
	}

	return
}

// -----------------------------------------------------------------------
// EqualsFloat32 checker.

var NoEqualsFloat32MoreThanMaxFloat32Error = "Comparing incomparable values as float32: one of parameters is more than math.MaxFloat32 / 2"
var NoEqualsFloat32LessThanMaxFloat32Error = "Comparing incomparable values as float32: one of parameters is less than -1 * math.MaxFloat32 / 2"

type equalsFloat32 struct {
	*CheckerInfo
}

// The EqualsFloat32 checker verifies that the obtained value is equal to
// the expected value as Equals checker, but ALWAYS converts parameters to float32 if it is possible.
//
// For example:
//
//     c.Assert(value, Equals, 42)
//     c.Assert(42.0, Equals, int64(42)) // true
//     c.Assert(int32(42), Equals, int64(42)) // true
// //
//    UnsupportedTypes:
//    []byte, string
//

var EqualsFloat32 Checker = &equalsFloat32{
	&CheckerInfo{Name: "EqualsFloat32", Params: []string{"get", "should be equal"}},
}

func (checker *equalsFloat32) Check(params []interface{}, names []string) (result bool, error string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			error = fmt.Sprint(v)
		} else if !result && error == "" {
			error = fmt.Sprintf("Difference: %v != %v", params[0], params[1])
		}
	}()

	// unsupported types
	error = "Comparing incomparable type as float32: " +
		reflect.ValueOf(params[0]).Type().String() +
		" and " +
		reflect.ValueOf(params[1]).Type().String()

	if isStringType(params[0]) || isStringType(params[1]) {
		return
	}

	rt := reflect.ValueOf(params[0]).Kind()
	if rt == reflect.Ptr {
		rt = reflect.ValueOf(params[0]).Type().Kind()
	}

	rt2 := reflect.ValueOf(params[1]).Kind()
	if rt2 == reflect.Ptr {
		rt2 = reflect.ValueOf(params[1]).Type().Kind()
	}

	switch rt {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		switch rt2 {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:

			if cf.Float32(params[0]) > math.MaxFloat32/2 || cf.Float32(params[1]) > math.MaxFloat32/2 {
				error = NoEqualsFloat32MoreThanMaxFloat32Error
			} else if cf.Float32(params[0]) < -1*math.MaxFloat32/2 || cf.Float32(params[1]) < -1*math.MaxFloat32/2 {
				error = NoEqualsFloat32LessThanMaxFloat32Error
			} else {
				result = cf.Float32(params[0]) == cf.Float32(params[1])
				error = ""
			}
		}
	}

	return
}
