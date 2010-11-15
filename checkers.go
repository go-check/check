/*
Gocheck - A rich testing framework for Go

Copyright (c) 2010, Gustavo Niemeyer <gustavo@niemeyer.net>

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

    * Redistributions of source code must retain the above copyright notice,
      this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright notice,
      this list of conditions and the following disclaimer in the documentation
      and/or other materials provided with the distribution.
    * Neither the name of the copyright holder nor the names of its
      contributors may be used to endorse or promote products derived from
      this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package gocheck

import (
    "reflect"
    "unsafe"
    "regexp"
    "fmt"
)


// -----------------------------------------------------------------------
// BugInfo and Bug() helper, to attach extra information to checks.

type bugInfo struct {
    format string
    args []interface{}
}

// Function to attach some information to an Assert() or Check() call.
// This should be used as, for instance:
//
//  Assert(a, Equals, 8192, Bug("Buffer size is incorrect, bug #123"))
//
// If the matching fails, the provided arguments will be passed to
// fmt.Sprintf(), and will be presented next to the logged failure. Note
// that it must be the last argument provided.
func Bug(format string, args ...interface{}) BugInfo {
    return &bugInfo{format, args}
}

// Interface which must be supported for attaching extra information to
// checks.  See the Bug() function.
type BugInfo interface {
    GetBugInfo() string
}

func (bug *bugInfo) GetBugInfo() string {
    return fmt.Sprintf(bug.format, bug.args...)
}


// -----------------------------------------------------------------------
// A useful Checker template.

// Checkers used with the c.Assert() and c.Check() verification methods
// must have this interface.  See the CheckerType type for an understanding
// of how the individual methods must work.
type Checker interface {
    Name() string
    VarNames() (obtained, expected string)
    NeedsExpectedValue() bool
    Check(obtained, expected interface{}) (result bool, error string)
}

// Sample checker type with some sane defaults.
type CheckerType struct{}

// Trick to ensure it matchers the desired interface.
var _ Checker = (*CheckerType)(nil)


// The function name used to build the matcher. E.g. "IsNil".
func (checker *CheckerType) Name() string {
    return "Checker"
}

// Method must return true if the given matcher needs to be informed
// of an expected value in addition to the actual value obtained to
// verify its expectations. E.g. false for IsNil.
func (checker *CheckerType) NeedsExpectedValue() bool {
    return true
}

// Variable names to be used for the obtained and expected values when
// reporting a failure in the expectations established. E.g.
// "obtained" and "expected".
func (checker *CheckerType) VarNames() (obtained, expected string) {
    return "obtained", "expected"
}

// Method must return true if the obtained value succeeds the
// expectations established by the given matcher.  If an error is
// returned, it means the provided parameters are somehow invalid.
func (checker *CheckerType) Check(obtained, expected interface{}) (
        result bool, error string) {
    return false, ""
}


// -----------------------------------------------------------------------
// Not() checker logic inverter.

// Invert the logic of the provided checker.  The resulting checker will
// succeed where the original one failed, and vice versa.  E.g.
// Assert(a, Not(Equals), b)
func Not(checker Checker) Checker {
    return &notChecker{checker}
}

type notChecker struct {
    sub Checker
}

func (checker *notChecker) Name() string {
    return "Not(" + checker.sub.Name() + ")"
}

func (checker *notChecker) NeedsExpectedValue() bool {
    return checker.sub.NeedsExpectedValue()
}

func (checker *notChecker) VarNames() (obtained, expected string) {
    obtained, expected = checker.sub.VarNames()
    return
}

func (checker *notChecker) Check(obtained, expected interface{}) (
        result bool, error string) {
    result, error = checker.sub.Check(obtained, expected)
    result = !result // So much for so little. :-)
    return
}


// -----------------------------------------------------------------------
// IsNil checker.

// Check if the obtained value is nil. E.g. Assert(err, IsNil).
var IsNil Checker = &isNilChecker{}

type isNilChecker struct{
    CheckerType
}

func (checker *isNilChecker) Name() string {
    return "IsNil"
}

func (checker *isNilChecker) NeedsExpectedValue() bool {
    return false
}

func (checker *isNilChecker) VarNames() (obtained, expected string) {
    return "value", ""
}

func (checker *isNilChecker) Check(obtained, expected interface{}) (
        result bool, error string) {
    return isNil(obtained), ""
}


func isNil(obtained interface{}) (result bool) {
    if obtained == nil {
        result = true
    } else {
        value := reflect.NewValue(obtained)
        result = *(*uintptr)(unsafe.Pointer(value.Addr())) == 0
    }
    return
}


// -----------------------------------------------------------------------
// NotNil checker. Alias for Not(IsNil), since it's so common.

// Check if the obtained value is not nil. E.g. Assert(iface, NotNil).
// This is an Alias for Not(IsNil), since it's a fairly common check.
var NotNil Checker = &notNilChecker{}

type notNilChecker struct{
    isNilChecker
}

func (checker *notNilChecker) Name() string {
    return "NotNil"
}

func (checker *notNilChecker) Check(obtained, expected interface{}) (
        result bool, error string) {
    return !isNil(obtained), ""
}


// -----------------------------------------------------------------------
// Equals checker.

// Check that the obtained value is equal to the expected value.  The
// check will work correctly even when facing arrays, interfaces, and
// values of different types (which always fails the test). E.g.
// Assert(value, Equals, 42).
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

// Check that the string provided as the obtained value (or the result of
// its .String() method, in case the value is not a string) matches the
// regular expression provided.  Note that, given the interface os.Error
// commonly used for errors, this checker will correctly verify its
// string representation. E.g. Assert(err, Matches, "perm.*denied")
var Matches Checker = &matchesChecker{}

type matchesChecker struct {
    CheckerType
}

func (checker *matchesChecker) Name() string {
    return "Matches"
}

func (checker *matchesChecker) VarNames() (obtained, expected string) {
    return "value", "regex"
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
