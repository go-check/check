package check

import (
	"errors"
	"reflect"
	"runtime"
)

type CheckersS struct{}

var _ = Suite(&CheckersS{})

func testInfo(c *C, checker Checker, name string, paramNames []string) {
	info := checker.Info()
	if info.Name != name {
		c.Fatalf("Got name %s, expected %s", info.Name, name)
	}
	if !reflect.DeepEqual(info.Params, paramNames) {
		c.Fatalf("Got param names %#v, expected %#v", info.Params, paramNames)
	}
}

func testCheck(c *C, checker Checker, result bool, error string, params ...interface{}) ([]interface{}, []string) {
	info := checker.Info()
	if len(params) != len(info.Params) {
		c.Fatalf("unexpected param count in test; expected %d got %d", len(info.Params), len(params))
	}
	names := append([]string{}, info.Params...)
	result_, error_ := checker.Check(params, names)
	if result_ != result || error_ != error {
		c.Fatalf("%s.Check(%#v) returned (%#v, %#v) rather than (%#v, %#v)",
			info.Name, params, result_, error_, result, error)
	}
	return params, names
}

func (s *CheckersS) TestComment(c *C) {
	bug := Commentf("a %d bc", 42)
	comment := bug.CheckCommentString()
	if comment != "a 42 bc" {
		c.Fatalf("Commentf returned %#v", comment)
	}
}

func (s *CheckersS) TestIsNil(c *C) {
	testInfo(c, IsNil, "IsNil", []string{"value"})

	testCheck(c, IsNil, true, "", nil)
	testCheck(c, IsNil, false, "", "a")

	testCheck(c, IsNil, true, "", (chan int)(nil))
	testCheck(c, IsNil, false, "", make(chan int))
	testCheck(c, IsNil, true, "", (error)(nil))
	testCheck(c, IsNil, false, "", errors.New(""))
	testCheck(c, IsNil, true, "", ([]int)(nil))
	testCheck(c, IsNil, false, "", make([]int, 1))
	testCheck(c, IsNil, false, "", int(0))
}

func (s *CheckersS) TestNotNil(c *C) {
	testInfo(c, NotNil, "NotNil", []string{"value"})

	testCheck(c, NotNil, false, "", nil)
	testCheck(c, NotNil, true, "", "a")

	testCheck(c, NotNil, false, "", (chan int)(nil))
	testCheck(c, NotNil, true, "", make(chan int))
	testCheck(c, NotNil, false, "", (error)(nil))
	testCheck(c, NotNil, true, "", errors.New(""))
	testCheck(c, NotNil, false, "", ([]int)(nil))
	testCheck(c, NotNil, true, "", make([]int, 1))
}

func (s *CheckersS) TestNot(c *C) {
	testInfo(c, Not(IsNil), "Not(IsNil)", []string{"value"})

	testCheck(c, Not(IsNil), false, "", nil)
	testCheck(c, Not(IsNil), true, "", "a")
}

type simpleStruct struct {
	i int
}

func (s *CheckersS) TestEquals(c *C) {
	testInfo(c, Equals, "Equals", []string{"obtained", "expected"})

	// The simplest.
	testCheck(c, Equals, true, "", 42, 42)
	testCheck(c, Equals, false, "", 42, 43)

	// Different native types.
	testCheck(c, Equals, false, "", int32(42), int64(42))

	// With nil.
	testCheck(c, Equals, false, "", 42, nil)

	// Slices
	testCheck(c, Equals, false, "runtime error: comparing uncomparable type []uint8", []byte{1, 2}, []byte{1, 2})

	// Struct values
	testCheck(c, Equals, true, "", simpleStruct{1}, simpleStruct{1})
	testCheck(c, Equals, false, "", simpleStruct{1}, simpleStruct{2})

	// Struct pointers
	testCheck(c, Equals, false, "", &simpleStruct{1}, &simpleStruct{1})
	testCheck(c, Equals, false, "", &simpleStruct{1}, &simpleStruct{2})
}

func (s *CheckersS) TestDeepEquals(c *C) {
	testInfo(c, DeepEquals, "DeepEquals", []string{"obtained", "expected"})

	// The simplest.
	testCheck(c, DeepEquals, true, "", 42, 42)
	testCheck(c, DeepEquals, false, "", 42, 43)

	// Different native types.
	testCheck(c, DeepEquals, false, "", int32(42), int64(42))

	// With nil.
	testCheck(c, DeepEquals, false, "", 42, nil)

	// Slices
	testCheck(c, DeepEquals, true, "", []byte{1, 2}, []byte{1, 2})
	testCheck(c, DeepEquals, false, "", []byte{1, 2}, []byte{1, 3})

	// Struct values
	testCheck(c, DeepEquals, true, "", simpleStruct{1}, simpleStruct{1})
	testCheck(c, DeepEquals, false, "", simpleStruct{1}, simpleStruct{2})

	// Struct pointers
	testCheck(c, DeepEquals, true, "", &simpleStruct{1}, &simpleStruct{1})
	testCheck(c, DeepEquals, false, "", &simpleStruct{1}, &simpleStruct{2})
}

func (s *CheckersS) TestHasLen(c *C) {
	testInfo(c, HasLen, "HasLen", []string{"obtained", "n"})

	testCheck(c, HasLen, true, "", "abcd", 4)
	testCheck(c, HasLen, true, "", []int{1, 2}, 2)
	testCheck(c, HasLen, false, "", []int{1, 2}, 3)

	testCheck(c, HasLen, false, "n must be an int", []int{1, 2}, "2")
	testCheck(c, HasLen, false, "obtained value type has no length", nil, 2)
}

func (s *CheckersS) TestErrorMatches(c *C) {
	testInfo(c, ErrorMatches, "ErrorMatches", []string{"value", "regex"})

	testCheck(c, ErrorMatches, false, "Error value is nil", nil, "some error")
	testCheck(c, ErrorMatches, false, "Value is not an error", 1, "some error")
	testCheck(c, ErrorMatches, true, "", errors.New("some error"), "some error")
	testCheck(c, ErrorMatches, true, "", errors.New("some error"), "so.*or")

	// Verify params mutation
	params, names := testCheck(c, ErrorMatches, false, "", errors.New("some error"), "other error")
	c.Assert(params[0], Equals, "some error")
	c.Assert(names[0], Equals, "error")
}

func (s *CheckersS) TestMatches(c *C) {
	testInfo(c, Matches, "Matches", []string{"value", "regex"})

	// Simple matching
	testCheck(c, Matches, true, "", "abc", "abc")
	testCheck(c, Matches, true, "", "abc", "a.c")

	// Must match fully
	testCheck(c, Matches, false, "", "abc", "ab")
	testCheck(c, Matches, false, "", "abc", "bc")

	// String()-enabled values accepted
	testCheck(c, Matches, true, "", reflect.ValueOf("abc"), "a.c")
	testCheck(c, Matches, false, "", reflect.ValueOf("abc"), "a.d")

	// Some error conditions.
	testCheck(c, Matches, false, "Obtained value is not a string and has no .String()", 1, "a.c")
	testCheck(c, Matches, false, "Can't compile regex: error parsing regexp: missing closing ]: `[c$`", "abc", "a[c")
}

func (s *CheckersS) TestPanics(c *C) {
	testInfo(c, Panics, "Panics", []string{"function", "expected"})

	// Some errors.
	testCheck(c, Panics, false, "Function has not panicked", func() bool { return false }, "BOOM")
	testCheck(c, Panics, false, "Function must take zero arguments", 1, "BOOM")

	// Plain strings.
	testCheck(c, Panics, true, "", func() { panic("BOOM") }, "BOOM")
	testCheck(c, Panics, false, "", func() { panic("KABOOM") }, "BOOM")
	testCheck(c, Panics, true, "", func() bool { panic("BOOM") }, "BOOM")

	// Error values.
	testCheck(c, Panics, true, "", func() { panic(errors.New("BOOM")) }, errors.New("BOOM"))
	testCheck(c, Panics, false, "", func() { panic(errors.New("KABOOM")) }, errors.New("BOOM"))

	type deep struct{ i int }
	// Deep value
	testCheck(c, Panics, true, "", func() { panic(&deep{99}) }, &deep{99})

	// Verify params/names mutation
	params, names := testCheck(c, Panics, false, "", func() { panic(errors.New("KABOOM")) }, errors.New("BOOM"))
	c.Assert(params[0], ErrorMatches, "KABOOM")
	c.Assert(names[0], Equals, "panic")

	// Verify a nil panic
	testCheck(c, Panics, true, "", func() { panic(nil) }, nil)
	testCheck(c, Panics, false, "", func() { panic(nil) }, "NOPE")
}

func (s *CheckersS) TestPanicMatches(c *C) {
	testInfo(c, PanicMatches, "PanicMatches", []string{"function", "expected"})

	// Error matching.
	testCheck(c, PanicMatches, true, "", func() { panic(errors.New("BOOM")) }, "BO.M")
	testCheck(c, PanicMatches, false, "", func() { panic(errors.New("KABOOM")) }, "BO.M")

	// Some errors.
	testCheck(c, PanicMatches, false, "Function has not panicked", func() bool { return false }, "BOOM")
	testCheck(c, PanicMatches, false, "Function must take zero arguments", 1, "BOOM")

	// Plain strings.
	testCheck(c, PanicMatches, true, "", func() { panic("BOOM") }, "BO.M")
	testCheck(c, PanicMatches, false, "", func() { panic("KABOOM") }, "BOOM")
	testCheck(c, PanicMatches, true, "", func() bool { panic("BOOM") }, "BO.M")

	// Verify params/names mutation
	params, names := testCheck(c, PanicMatches, false, "", func() { panic(errors.New("KABOOM")) }, "BOOM")
	c.Assert(params[0], Equals, "KABOOM")
	c.Assert(names[0], Equals, "panic")

	// Verify a nil panic
	testCheck(c, PanicMatches, false, "Panic value is not a string or an error", func() { panic(nil) }, "")
}

func (s *CheckersS) TestFitsTypeOf(c *C) {
	testInfo(c, FitsTypeOf, "FitsTypeOf", []string{"obtained", "sample"})

	// Basic types
	testCheck(c, FitsTypeOf, true, "", 1, 0)
	testCheck(c, FitsTypeOf, false, "", 1, int64(0))

	// Aliases
	testCheck(c, FitsTypeOf, false, "", 1, errors.New(""))
	testCheck(c, FitsTypeOf, false, "", "error", errors.New(""))
	testCheck(c, FitsTypeOf, true, "", errors.New("error"), errors.New(""))

	// Structures
	testCheck(c, FitsTypeOf, false, "", 1, simpleStruct{})
	testCheck(c, FitsTypeOf, false, "", simpleStruct{42}, &simpleStruct{})
	testCheck(c, FitsTypeOf, true, "", simpleStruct{42}, simpleStruct{})
	testCheck(c, FitsTypeOf, true, "", &simpleStruct{42}, &simpleStruct{})

	// Some bad values
	testCheck(c, FitsTypeOf, false, "Invalid sample value", 1, interface{}(nil))
	testCheck(c, FitsTypeOf, false, "", interface{}(nil), 0)
}

func (s *CheckersS) TestImplements(c *C) {
	testInfo(c, Implements, "Implements", []string{"obtained", "ifaceptr"})

	var e error
	var re runtime.Error
	testCheck(c, Implements, true, "", errors.New(""), &e)
	testCheck(c, Implements, false, "", errors.New(""), &re)

	// Some bad values
	testCheck(c, Implements, false, "ifaceptr should be a pointer to an interface variable", 0, errors.New(""))
	testCheck(c, Implements, false, "ifaceptr should be a pointer to an interface variable", 0, interface{}(nil))
	testCheck(c, Implements, false, "", interface{}(nil), &e)
}
