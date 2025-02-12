package check_test

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"

	"github.com/iostrovok/check"
)

type CheckersS struct{}

var _ = check.Suite(&CheckersS{})

func testInfo(c *check.C, checker check.Checker, name string, paramNames []string) {
	info := checker.Info()
	if info.Name != name {
		c.Fatalf("Got name %s, expected %s", info.Name, name)
	}

	if !reflect.DeepEqual(info.Params, paramNames) {
		c.Fatalf("Got param names %#v, expected %#v", info.Params, paramNames)
	}
}

func testCheck(c *check.C, checker check.Checker, result bool, error string, params ...any) ([]any, []string) {
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

func testCheckNoLine(c *check.C, checker check.Checker, result bool, error string, params ...any) ([]any, []string) {
	info := checker.Info()
	if len(params) != len(info.Params) {
		c.Fatalf("unexpected param count in test; expected %d got %d", len(info.Params), len(params))
	}
	names := append([]string{}, info.Params...)
	result_, error_ := checker.Check(params, names)
	if result_ != result {
		c.Fatalf("%s.Check(%#v) returned (%#v, %#v) rather than (%#v, %#v)",
			info.Name, params, result_, error_, result, error)
	}
	return params, names
}

func (s *CheckersS) TestComment(c *check.C) {
	bug := check.Commentf("a %d bc", 42)
	comment := bug.CheckCommentString()
	if comment != "a 42 bc" {
		c.Fatalf("Commentf returned %#v", comment)
	}
}

func (s *CheckersS) TestIsNil(c *check.C) {
	testInfo(c, check.IsNil, "IsNil", []string{"value"})

	testCheck(c, check.IsNil, true, "", nil)
	testCheck(c, check.IsNil, false, "", "a")

	testCheck(c, check.IsNil, true, "", (chan int)(nil))
	testCheck(c, check.IsNil, false, "", make(chan int))
	testCheck(c, check.IsNil, true, "", (error)(nil))
	testCheck(c, check.IsNil, false, "", errors.New(""))
	testCheck(c, check.IsNil, true, "", ([]int)(nil))
	testCheck(c, check.IsNil, false, "", make([]int, 1))
	testCheck(c, check.IsNil, false, "", int(0))
}

func (s *CheckersS) TestIsFalse(c *check.C) {
	testInfo(c, check.IsFalse, "IsFalse", []string{"value"})

	testCheck(c, check.IsFalse, true, "", nil)
	testCheck(c, check.IsFalse, false, "", "a")

	testCheck(c, check.IsFalse, false, "", (chan int)(nil))
	testCheck(c, check.IsFalse, false, "", make(chan int))
	testCheck(c, check.IsFalse, true, "", (error)(nil))
	testCheck(c, check.IsFalse, true, "", errors.New(""))
	testCheck(c, check.IsFalse, true, "", ([]int)(nil))
	testCheck(c, check.IsFalse, false, "", make([]int, 1))
	testCheck(c, check.IsFalse, true, "", 0)
}

func (s *CheckersS) TestNotEmpty(c *check.C) {
	testInfo(c, check.NotEmpty, "NotEmpty", []string{"value"})

	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '0'", 0)
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '0'", 0.0)
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '<nil>'", nil)

	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '<nil>'", (error)(nil))
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was ''", errors.New(""))
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '[]'", ([]int)(nil))
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '[]'", make([]int, 0))
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '[]'", []int{})
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was 'map[]'", map[int]map[string]bool{})

	testCheck(c, check.NotEmpty, true, "", map[int]bool{1: true})
	testCheck(c, check.NotEmpty, true, "", make([]int, 1))

	testCheck(c, check.NotEmpty, true, "", "a")
	testCheck(c, check.NotEmpty, true, "", 100)
	testCheck(c, check.NotEmpty, true, "", -1)

	// Channels!
	testCheck(c, check.NotEmpty, false, "should NOT be empty, but was '<nil>'", (chan int)(nil))

	ch := make(chan int)
	testCheck(c, check.NotEmpty, false, fmt.Sprintf("should NOT be empty, but was '%+v'", ch), ch)

	ch = make(chan int, 1)
	testCheck(c, check.NotEmpty, false, fmt.Sprintf("should NOT be empty, but was '%+v'", ch), ch)
	ch <- 1
	testCheck(c, check.NotEmpty, true, "", ch)
}

func (s *CheckersS) TestIsEmpty(c *check.C) {
	testInfo(c, check.IsEmpty, "IsEmpty", []string{"value"})

	testCheck(c, check.IsEmpty, true, "", 0)
	testCheck(c, check.IsEmpty, true, "", 0.0)
	testCheck(c, check.IsEmpty, true, "", nil)

	testCheck(c, check.IsEmpty, true, "", (error)(nil))
	testCheck(c, check.IsEmpty, true, "", errors.New(""))
	testCheck(c, check.IsEmpty, true, "", ([]int)(nil))
	testCheck(c, check.IsEmpty, true, "", make([]int, 0))
	testCheck(c, check.IsEmpty, true, "", []int{})
	testCheck(c, check.IsEmpty, true, "", map[int]map[string]bool{})

	testCheck(c, check.IsEmpty, false, "should be empty, but was 'map[1:true]'", map[int]bool{1: true})
	testCheck(c, check.IsEmpty, false, "should be empty, but was '[0]'", make([]int, 1))
	testCheck(c, check.IsEmpty, false, "should be empty, but was 'a'", "a")
	testCheck(c, check.IsEmpty, false, "should be empty, but was '100'", 100)
	testCheck(c, check.IsEmpty, false, "should be empty, but was '-1'", -1)

	// Channels!
	testCheck(c, check.IsEmpty, true, "", (chan int)(nil))
	testCheck(c, check.IsEmpty, true, "", make(chan int))

	ch := make(chan int)
	testCheck(c, check.IsEmpty, true, "", ch)

	ch = make(chan int, 1)
	testCheck(c, check.IsEmpty, true, "", ch)
	ch <- 1
	testCheck(c, check.IsEmpty, false, fmt.Sprintf("should be empty, but was '%+v'", ch), ch)
}

func (s *CheckersS) TestIsTrue(c *check.C) {
	testInfo(c, check.IsTrue, "IsTrue", []string{"value"})

	testCheck(c, check.IsTrue, false, "", nil)
	testCheck(c, check.IsTrue, true, "", "a")

	testCheck(c, check.IsTrue, true, "", (chan int)(nil))
	testCheck(c, check.IsTrue, true, "", make(chan int))
	testCheck(c, check.IsTrue, false, "", (error)(nil))
	testCheck(c, check.IsTrue, false, "", errors.New(""))
	testCheck(c, check.IsTrue, false, "", ([]int)(nil))
	testCheck(c, check.IsTrue, true, "", make([]int, 1))
	testCheck(c, check.IsTrue, false, "", 0)
}

func (s *CheckersS) TestNotNil(c *check.C) {
	testInfo(c, check.NotNil, "NotNil", []string{"value"})

	testCheck(c, check.NotNil, false, "", nil)
	testCheck(c, check.NotNil, true, "", "a")

	testCheck(c, check.NotNil, false, "", (chan int)(nil))
	testCheck(c, check.NotNil, true, "", make(chan int))
	testCheck(c, check.NotNil, false, "", (error)(nil))
	testCheck(c, check.NotNil, true, "", errors.New(""))
	testCheck(c, check.NotNil, false, "", ([]int)(nil))
	testCheck(c, check.NotNil, true, "", make([]int, 1))
}

func (s *CheckersS) TestNot(c *check.C) {
	testInfo(c, check.Not(check.IsNil), "Not(IsNil)", []string{"value"})

	testCheck(c, check.Not(check.IsNil), false, "", nil)
	testCheck(c, check.Not(check.IsNil), true, "", "a")
	testCheck(c, check.Not(check.Equals), true, "", 42, 43)
}

type simpleStruct struct {
	i int
}

func (s *CheckersS) TestEquals(c *check.C) {
	testInfo(c, check.Equals, "Equals", []string{"obtained", "expected"})

	// The simplest.
	testCheck(c, check.Equals, true, "", 42, 42)
	testCheck(c, check.Equals, false, "", 42, 43)

	// Different native types.
	testCheck(c, check.Equals, false, "", int32(42), int64(42))

	// With nil.
	testCheck(c, check.Equals, false, "", 42, nil)
	testCheck(c, check.Equals, false, "", nil, 42)
	testCheck(c, check.Equals, true, "", nil, nil)

	// Slices
	testCheck(c, check.Equals, false, "runtime error: comparing uncomparable type []uint8", []byte{1, 2}, []byte{1, 2})

	// Struct values
	testCheck(c, check.Equals, true, "", simpleStruct{1}, simpleStruct{1})
	testCheck(c, check.Equals, false, `Difference:
...     i: 1 != 2
`, simpleStruct{1}, simpleStruct{2})

	// Struct pointers, no difference in values, just pointer
	testCheck(c, check.Equals, false, "", &simpleStruct{1}, &simpleStruct{1})
	// Struct pointers, different pointers and different values
	testCheck(c, check.Equals, false, `Difference:
...     i: 1 != 2
`, &simpleStruct{1}, &simpleStruct{2})
}

func (s *CheckersS) TestDeepEquals(c *check.C) {
	testInfo(c, check.DeepEquals, "DeepEquals", []string{"obtained", "expected"})

	// The simplest.
	testCheck(c, check.DeepEquals, true, "", 42, 42)
	testCheck(c, check.DeepEquals, false, "", 42, 43)

	// Different native types.
	testCheck(c, check.DeepEquals, false, "", int32(42), int64(42))

	// With nil.
	testCheck(c, check.DeepEquals, false, "", 42, nil)

	// Slices
	testCheck(c, check.DeepEquals, true, "", []byte{1, 2}, []byte{1, 2})
	testCheck(c, check.DeepEquals, false, `Difference:
...     [1]: 2 != 3
`, []byte{1, 2}, []byte{1, 3})

	// Struct values
	testCheck(c, check.DeepEquals, true, "", simpleStruct{1}, simpleStruct{1})
	testCheck(c, check.DeepEquals, false, `Difference:
...     i: 1 != 2
`, simpleStruct{1}, simpleStruct{2})

	// Struct pointers
	testCheck(c, check.DeepEquals, true, "", &simpleStruct{1}, &simpleStruct{1})
	s1 := &simpleStruct{1}
	s2 := &simpleStruct{2}
	testCheck(c, check.DeepEquals, false, `Difference:
...     i: 1 != 2
`, s1, s2)
}

func (s *CheckersS) TestHasLen(c *check.C) {
	testInfo(c, check.HasLen, "HasLen", []string{"obtained", "n"})

	testCheck(c, check.HasLen, true, "", "abcd", 4)
	testCheck(c, check.HasLen, true, "", []int{1, 2}, 2)
	testCheck(c, check.HasLen, false, "", []int{1, 2}, 3)

	testCheck(c, check.HasLen, false, "n must be an int*, not string", []int{1, 2}, "2")
	testCheck(c, check.HasLen, false, "obtained value type has no length property", nil, 2)
}

func (s *CheckersS) TestHasLenLessThan(c *check.C) {
	testInfo(c, check.HasLen, "HasLen", []string{"obtained", "n"})

	testCheck(c, check.HasLenLessThan, true, "", "abcd", int8(32))
	testCheck(c, check.HasLenLessThan, true, "", "abcd", 5)
	testCheck(c, check.HasLenLessThan, true, "", []int{1, 2, 3}, 5)
	testCheck(c, check.HasLenLessThan, true, "", []string{"1", "2", "3", "4"}, 7)
	testCheck(c, check.HasLenLessThan, true, "", []string{}, int64(2))

	testCheck(c, check.HasLenLessThan, false, "", []int{1, 2}, 2)
	testCheck(c, check.HasLenLessThan, false, "", []int{1, 2}, 1)
	testCheck(c, check.HasLenLessThan, false, "n must be an int*, not string", []int{1, 2}, "2")
	testCheck(c, check.HasLenLessThan, false, "obtained value type has no length property", nil, 2)
}

func (s *CheckersS) TestHasLenMoreThan(c *check.C) {
	testInfo(c, check.HasLen, "HasLen", []string{"obtained", "n"})

	testCheck(c, check.HasLenMoreThan, true, "", "abcd", 3)
	testCheck(c, check.HasLenMoreThan, true, "", "abcd", 0)
	testCheck(c, check.HasLenMoreThan, true, "", []int{1, 2, 3}, 2)
	testCheck(c, check.HasLenMoreThan, true, "", []string{"1", "2", "3", "4"}, 2)
	testCheck(c, check.HasLenMoreThan, true, "", []string{}, int64(-2))

	testCheck(c, check.HasLenMoreThan, true, "", map[int]any{1: 1, 2: "", 3: 0, 4: "10"}, uint(3))

	testCheck(c, check.HasLenMoreThan, false, "", []int{1, 2}, 2)
	testCheck(c, check.HasLenMoreThan, false, "", []int{1, 2}, 3)
	testCheck(c, check.HasLenMoreThan, false, "n must be an *int, not string", []int{1, 2}, "2")
	testCheck(c, check.HasLenMoreThan, false, "obtained value type has no length property", nil, 2)
}

func (s *CheckersS) TestErrorMatches(c *check.C) {
	testInfo(c, check.ErrorMatches, "ErrorMatches", []string{"value", "regex"})

	testCheck(c, check.ErrorMatches, false, "Error value is nil", nil, "some error")
	testCheck(c, check.ErrorMatches, false, "Value is not an error", 1, "some error")
	testCheck(c, check.ErrorMatches, true, "", errors.New("some error"), "some error")
	testCheck(c, check.ErrorMatches, true, "", errors.New("some error"), "so.*or")

	// Verify params mutation
	params, names := testCheck(c, check.ErrorMatches, false, "", errors.New("some error"), "other error")
	c.Assert(params[0], check.Equals, "some error")
	c.Assert(names[0], check.Equals, "error")
}

func (s *CheckersS) TestMatches(c *check.C) {
	testInfo(c, check.Matches, "Matches", []string{"value", "regex"})

	// Simple matching
	testCheck(c, check.Matches, true, "", "abc", "abc")
	testCheck(c, check.Matches, true, "", "abc", "a.c")

	// Must match fully
	testCheck(c, check.Matches, false, "", "abc", "ab")
	testCheck(c, check.Matches, false, "", "abc", "bc")

	// String()-enabled values accepted
	testCheck(c, check.Matches, true, "", reflect.ValueOf("abc"), "a.c")
	testCheck(c, check.Matches, false, "", reflect.ValueOf("abc"), "a.d")

	// Some error conditions.
	testCheck(c, check.Matches, false, "Obtained value is not a string and has no .String()", 1, "a.c")
	testCheck(c, check.Matches, false, "Can't compile regex: error parsing regexp: missing closing ]: `[c$`", "abc", "a[c")
}

func (s *CheckersS) TestPanics(c *check.C) {
	testInfo(c, check.Panics, "Panics", []string{"function", "expected"})

	// Some errors.
	testCheck(c, check.Panics, false, "Function has not panicked", func() bool { return false }, "BOOM")
	testCheck(c, check.Panics, false, "Function must take zero arguments", 1, "BOOM")

	// Plain strings.
	testCheck(c, check.Panics, true, "", func() { panic("BOOM") }, "BOOM")
	testCheck(c, check.Panics, false, "", func() { panic("KABOOM") }, "BOOM")
	testCheck(c, check.Panics, true, "", func() bool { panic("BOOM") }, "BOOM")

	// Error values.
	testCheck(c, check.Panics, true, "", func() { panic(errors.New("BOOM")) }, errors.New("BOOM"))
	testCheck(c, check.Panics, false, "", func() { panic(errors.New("KABOOM")) }, errors.New("BOOM"))

	type deep struct{ i int }
	// Deep value
	testCheck(c, check.Panics, true, "", func() { panic(&deep{99}) }, &deep{99})

	// Verify params/names mutation
	params, names := testCheck(c, check.Panics, false, "", func() { panic(errors.New("KABOOM")) }, errors.New("BOOM"))
	c.Assert(params[0], check.ErrorMatches, "KABOOM")
	c.Assert(names[0], check.Equals, "panic")

	// Verify a nil panic
	testCheck(c, check.Panics, true, "", func() { panic(nil) }, nil)
	testCheck(c, check.Panics, false, "", func() { panic(nil) }, "NOPE")
}

func (s *CheckersS) TestPanicMatches(c *check.C) {
	testInfo(c, check.PanicMatches, "PanicMatches", []string{"function", "expected"})

	// Error matching.
	testCheck(c, check.PanicMatches, true, "", func() { panic(errors.New("BOOM")) }, "BO.M")
	testCheck(c, check.PanicMatches, false, "", func() { panic(errors.New("KABOOM")) }, "BO.M")

	// Some errors.
	testCheck(c, check.PanicMatches, false, "Function has not panicked", func() bool { return false }, "BOOM")
	testCheck(c, check.PanicMatches, false, "Function must take zero arguments", 1, "BOOM")

	// Plain strings.
	testCheck(c, check.PanicMatches, true, "", func() { panic("BOOM") }, "BO.M")
	testCheck(c, check.PanicMatches, false, "", func() { panic("KABOOM") }, "BOOM")
	testCheck(c, check.PanicMatches, true, "", func() bool { panic("BOOM") }, "BO.M")

	// Verify params/names mutation
	params, names := testCheck(c, check.PanicMatches, false, "", func() { panic(errors.New("KABOOM")) }, "BOOM")
	c.Assert(params[0], check.Equals, "KABOOM")
	c.Assert(names[0], check.Equals, "panic")
}

func (s *CheckersS) TestPanicMatchesNil(c *check.C) {
	c.Skip("PanicMatches doesn't work with nil panics")
	// Verify a nil panic
	testCheck(c, check.PanicMatches, false, "panic called with nil argument", func() { panic(nil) }, "")
}

func (s *CheckersS) TestFitsTypeOf(c *check.C) {
	testInfo(c, check.FitsTypeOf, "FitsTypeOf", []string{"obtained", "sample"})

	// Basic types
	testCheck(c, check.FitsTypeOf, true, "", 1, 0)
	testCheck(c, check.FitsTypeOf, false, "", 1, int64(0))

	// Aliases
	testCheck(c, check.FitsTypeOf, false, "", 1, errors.New(""))
	testCheck(c, check.FitsTypeOf, false, "", "error", errors.New(""))
	testCheck(c, check.FitsTypeOf, true, "", errors.New("error"), errors.New(""))

	// Structures
	testCheck(c, check.FitsTypeOf, false, "", 1, simpleStruct{})
	testCheck(c, check.FitsTypeOf, false, "", simpleStruct{42}, &simpleStruct{})
	testCheck(c, check.FitsTypeOf, true, "", simpleStruct{42}, simpleStruct{})
	testCheck(c, check.FitsTypeOf, true, "", &simpleStruct{42}, &simpleStruct{})

	// Some bad values
	testCheck(c, check.FitsTypeOf, false, "Invalid sample value", 1, any(nil))
	testCheck(c, check.FitsTypeOf, false, "", any(nil), 0)
}

func (s *CheckersS) TestImplements(c *check.C) {
	testInfo(c, check.Implements, "Implements", []string{"obtained", "ifaceptr"})

	var e error
	var re runtime.Error
	testCheck(c, check.Implements, true, "", errors.New(""), &e)
	testCheck(c, check.Implements, false, "", errors.New(""), &re)

	// Some bad values
	testCheck(c, check.Implements, false, "ifaceptr should be a pointer to an interface variable", 0, errors.New(""))
	testCheck(c, check.Implements, false, "ifaceptr should be a pointer to an interface variable", 0, any(nil))
	testCheck(c, check.Implements, false, "", any(nil), &e)
}

func (s *CheckersS) TestMoreThan(c *check.C) {
	testCheck(c, check.MoreThan, true, "", 43, 42)
	testCheck(c, check.MoreThan, true, "", float32(43.11), float32(43.1))
	testCheck(c, check.MoreThan, true, "", 43.1201, 43.12)
	testCheck(c, check.MoreThan, true, "", uint64(4554), uint64(4455))
	testCheck(c, check.MoreThan, true, "", uint64(333444), uint64(111244))
	testCheck(c, check.MoreThan, false, "Comparing incomparable type int and float64", 43342, math.MaxFloat32+1000)

	testCheck(c, check.MoreThan, false, "Comparing incomparable type float32 and uint64", float32(43.12), uint64(43))
	testCheck(c, check.MoreThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32+1000, float32(math.MaxFloat32-1000))
	testCheck(c, check.MoreThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32, float32(math.MaxFloat32+1))
	testCheckNoLine(c, check.MoreThan, false, "Difference: 3.4028234663852886e+38 >= 3.4028234663852886e+38", math.MaxFloat32, math.MaxFloat32)
	testCheckNoLine(c, check.MoreThan, false, "Difference: -3.4028234663852886e+38 >= -3.4028234663852886e+38", -1*math.MaxFloat32, -1*math.MaxFloat32)

	// strings
	testCheck(c, check.MoreThan, true, "", "43", "42")
	testCheck(c, check.MoreThan, false, "First (string) parameter equals (string) second, expect more", []byte("42"), "42")
	testCheck(c, check.MoreThan, true, "", "423", []byte("421"))
	testCheck(c, check.MoreThan, false, check.NoLessThanStringError, []byte("42"), "421")

	testCheck(c, check.MoreThan, false, check.NoLessThanStringError, "Abc", []byte("abc"))
	testCheck(c, check.MoreThan, false, "Comparing incomparable type int and string", 41, "42")
	testCheck(c, check.MoreThan, false, check.NoLessThanStringError, "ABC", "abc")
	testCheck(c, check.MoreThan, false, "Comparing incomparable type []uint8 and float64", []byte("42"), float64(123.543))
}

func (s *CheckersS) TestLessThan(c *check.C) {
	testCheck(c, check.LessThan, true, "", 42, 43)
	testCheck(c, check.LessThan, true, "", float32(43.1), float32(43.12))
	testCheck(c, check.LessThan, true, "", 43.12, 43.1201)
	testCheck(c, check.LessThan, true, "", uint64(4454), uint64(4555))
	testCheck(c, check.LessThan, true, "", uint64(111244), uint64(333444))
	testCheck(c, check.LessThan, false, "Comparing incomparable type int and float64", 43342, math.MaxFloat32+1000)

	testCheck(c, check.LessThan, false, "Comparing incomparable type float32 and uint64", float32(43.12), uint64(43))
	testCheck(c, check.LessThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32+1000, float32(math.MaxFloat32-1000))
	testCheck(c, check.LessThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32, float32(math.MaxFloat32+1))
	testCheckNoLine(c, check.LessThan, false, "Difference: 3.4028234663852886e+38 >= 3.4028234663852886e+38", math.MaxFloat32, math.MaxFloat32)
	testCheckNoLine(c, check.LessThan, false, "Difference: -3.4028234663852886e+38 >= -3.4028234663852886e+38", -1*math.MaxFloat32, -1*math.MaxFloat32)

	// strings
	testCheck(c, check.LessThan, true, "", "42", "43")
	testCheck(c, check.LessThan, false, "First (string) parameter equals (string) second, expect less", []byte("43"), "43")
	testCheck(c, check.LessThan, true, "", []byte("42"), "43")
	testCheck(c, check.LessThan, true, "", "421", []byte("423"))
	testCheck(c, check.LessThan, false, check.NoLessThanStringError, []byte("421"), "42")

	testCheck(c, check.LessThan, false, check.NoLessThanStringError, "abc", []byte("Abc"))
	testCheck(c, check.LessThan, false, "Comparing incomparable type int and string", 41, "42")
	testCheck(c, check.LessThan, false, check.NoLessThanStringError, "abc", "ABC")
	testCheck(c, check.LessThan, false, "Comparing incomparable type []uint8 and float64", []byte("42"), float64(123.543))
}

func (s *CheckersS) TestMoreOrEqualThan(c *check.C) {
	testCheck(c, check.MoreOrEqualThan, true, "", 43, 42)
	testCheck(c, check.MoreOrEqualThan, true, "", float32(43.12), float32(43.12))
	testCheck(c, check.MoreOrEqualThan, true, "", 43.1201, 43.12)
	testCheck(c, check.MoreOrEqualThan, true, "", uint64(44), uint64(44))
	testCheck(c, check.MoreOrEqualThan, true, "", uint64(333444), uint64(111244))
	testCheck(c, check.MoreOrEqualThan, false, "Comparing incomparable type int and float64", 43342, math.MaxFloat32+1000)

	//
	testCheck(c, check.MoreOrEqualThan, false, "Comparing incomparable type float32 and uint64", float32(43.12), uint64(43))
	testCheck(c, check.MoreOrEqualThan, false, "Difference: 42 < 43", 42, 43)
	testCheck(c, check.MoreOrEqualThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32+1000, float32(math.MaxFloat32-1000))
	testCheck(c, check.MoreOrEqualThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32, float32(math.MaxFloat32-1))
	testCheck(c, check.MoreOrEqualThan, true, "", math.MaxFloat32, math.MaxFloat32)
	testCheck(c, check.MoreOrEqualThan, true, "", -1*math.MaxFloat32, -1*math.MaxFloat32)

	// strings
	testCheck(c, check.MoreOrEqualThan, true, "", "42", "42")
	testCheck(c, check.MoreOrEqualThan, true, "", []byte("42"), "42")
	testCheck(c, check.MoreOrEqualThan, true, "", []byte("421"), "42")
	testCheck(c, check.MoreOrEqualThan, false, check.NoMoreOrEqualThanStringError, []byte("42"), "421")
	//
	testCheck(c, check.MoreOrEqualThan, false, check.NoMoreOrEqualThanStringError, []byte("42"), "421")
	testCheck(c, check.MoreOrEqualThan, false, check.NoMoreOrEqualThanStringError, "Abc", []byte("abc"))
	testCheck(c, check.MoreOrEqualThan, false, "Comparing incomparable type int and string", 41, "42")
	testCheck(c, check.MoreOrEqualThan, false, check.NoMoreOrEqualThanStringError, "ABC", "abc")
	testCheck(c, check.MoreOrEqualThan, false, "Comparing incomparable type []uint8 and float64", []byte("42"), float64(123.543))
	testCheck(c, check.MoreOrEqualThan, false, check.NoMoreOrEqualThanStringError, "42", []byte("421"))
}

func (s *CheckersS) TestLessOrEqualThan(c *check.C) {
	testCheck(c, check.LessOrEqualThan, true, "", 42, 43)
	testCheck(c, check.LessOrEqualThan, true, "", float32(43.12), float32(43.12))
	testCheck(c, check.LessOrEqualThan, true, "", float64(43.12), float64(43.1201))
	testCheck(c, check.LessOrEqualThan, true, "", uint64(44), uint64(44))
	testCheck(c, check.LessOrEqualThan, true, "", uint64(111244), uint64(333444))
	testCheck(c, check.LessOrEqualThan, false, "Comparing incomparable type int and float64", 43342, math.MaxFloat32+1000)

	//
	testCheck(c, check.LessOrEqualThan, false, "Comparing incomparable type float32 and uint64", float32(43.12), uint64(43))
	testCheck(c, check.LessOrEqualThan, false, "Difference: 43 > 42", 43, 42)
	testCheck(c, check.LessOrEqualThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32+1000, float32(math.MaxFloat32-1000))
	testCheck(c, check.LessOrEqualThan, false, "Comparing incomparable type float64 and float32", math.MaxFloat32, float32(math.MaxFloat32-1))
	testCheck(c, check.LessOrEqualThan, true, "", math.MaxFloat32, math.MaxFloat32)
	testCheck(c, check.LessOrEqualThan, true, "", -1*math.MaxFloat32, -1*math.MaxFloat32)

	// strings
	testCheck(c, check.LessOrEqualThan, true, "", "42", "42")
	testCheck(c, check.LessOrEqualThan, true, "", []byte("42"), "42")
	testCheck(c, check.LessOrEqualThan, true, "", []byte("42"), "421")
	testCheck(c, check.LessOrEqualThan, false, check.NoLessOrEqualThanStringError, []byte("421"), "42")
	//
	testCheck(c, check.LessOrEqualThan, false, check.NoLessOrEqualThanStringError, []byte("421"), "42")
	testCheck(c, check.LessOrEqualThan, false, check.NoLessOrEqualThanStringError, "abc", []byte("Abc"))
	testCheck(c, check.LessOrEqualThan, false, "Comparing incomparable type int and string", 41, "42")
	testCheck(c, check.LessOrEqualThan, false, check.NoLessOrEqualThanStringError, "abc", "ABC")
	testCheck(c, check.LessOrEqualThan, false, "Comparing incomparable type []uint8 and float64", []byte("42"), float64(123.543))
	testCheck(c, check.LessOrEqualThan, false, check.NoLessOrEqualThanStringError, "421", []byte("42"))
}

func (s *CheckersS) TestEqualsMore(c *check.C) {
	testCheck(c, check.EqualsMore, true, "", 42, 42)
	testCheck(c, check.EqualsMore, true, "", float32(43.12), float32(43.12))
	testCheck(c, check.EqualsMore, true, "", float64(43.12), float64(43.12))
	testCheck(c, check.EqualsMore, true, "", uint64(44), uint64(44))
	testCheck(c, check.EqualsMore, true, "", uint64(111244), uint64(111244))
	testCheck(c, check.EqualsMore, false, "Comparing incomparable type int and float64", 43342, math.MaxFloat32+1000)

	//
	testCheck(c, check.EqualsMore, false, "Comparing incomparable type float32 and uint64", float32(43), uint64(43))
	testCheck(c, check.EqualsMore, false, "Difference: 43 != 42", 43, 42)
	testCheck(c, check.EqualsMore, false, "Difference: 43 != 64", 43, int32(64))
	testCheck(c, check.EqualsMore, false, "Comparing incomparable type float64 and float32", math.MaxFloat32+1000, float32(math.MaxFloat32-1000))
	testCheck(c, check.EqualsMore, false, "Comparing incomparable type float64 and float32", math.MaxFloat32, float32(math.MaxFloat32-1))
	testCheck(c, check.EqualsMore, true, "", math.MaxFloat32, math.MaxFloat32)
	testCheck(c, check.EqualsMore, true, "", -1*math.MaxFloat32, -1*math.MaxFloat32)

	// strings
	testCheck(c, check.EqualsMore, true, "", "42", "42")
	testCheck(c, check.EqualsMore, true, "", "42", []byte("42"))
	testCheck(c, check.EqualsMore, true, "", []byte("421"), "421")
	//
	testCheck(c, check.EqualsMore, false, check.NoEqualsMoreStringError, []byte("421"), "42")
	testCheck(c, check.EqualsMore, false, check.NoEqualsMoreStringError, "421", "42")
	testCheck(c, check.EqualsMore, false, check.NoEqualsMoreStringError, []byte("421"), "42")
	testCheck(c, check.EqualsMore, false, check.NoEqualsMoreStringError, "abc", []byte("Abc"))
	testCheck(c, check.EqualsMore, false, "Comparing incomparable type int and string", 41, "42")
	testCheck(c, check.EqualsMore, false, check.NoEqualsMoreStringError, "abc", "ABC")
	testCheck(c, check.EqualsMore, false, "Comparing incomparable type []uint8 and float64", []byte("42"), float64(123.543))
	testCheck(c, check.EqualsMore, false, check.NoEqualsMoreStringError, "421", []byte("42"))
}

func (s *CheckersS) TestEqualsFloat32(c *check.C) {
	testCheck(c, check.EqualsFloat32, true, "", 42, 42)
	testCheck(c, check.EqualsFloat32, true, "", float32(43.12), float32(43.12))
	testCheck(c, check.EqualsFloat32, true, "", float64(43.12), float64(43.12))
	testCheck(c, check.EqualsFloat32, true, "", uint64(44), uint64(44))
	testCheck(c, check.EqualsFloat32, true, "", uint64(111244), uint64(111244))
	testCheckNoLine(c, check.EqualsFloat32, false, "", 43342, math.MaxFloat32+1000)

	//
	testCheck(c, check.EqualsFloat32, true, "", float32(43), uint64(43))
	testCheck(c, check.EqualsFloat32, false, "Difference: 43 != 42", 43, 42)
	testCheck(c, check.EqualsFloat32, false, "Difference: 43 != 64", 43, int32(64))
	testCheck(c, check.EqualsFloat32, false, check.NoEqualsFloat32MoreThanMaxFloat32Error, math.MaxFloat32+1000, float32(math.MaxFloat32-1000))
	testCheck(c, check.EqualsFloat32, false, check.NoEqualsFloat32MoreThanMaxFloat32Error, math.MaxFloat32, float32(math.MaxFloat32-1))
	testCheck(c, check.EqualsFloat32, false, check.NoEqualsFloat32MoreThanMaxFloat32Error, math.MaxFloat32, math.MaxFloat32)
	testCheck(c, check.EqualsFloat32, false, check.NoEqualsFloat32LessThanMaxFloat32Error, -1*math.MaxFloat32, -1*math.MaxFloat32)

	// strings
	testCheck(c, check.EqualsFloat32, false, "Comparing incomparable type as float32: string and int", "42", 42)
	testCheck(c, check.EqualsFloat32, false, "Comparing incomparable type as float32: string and string", "42", "42")
	testCheck(c, check.EqualsFloat32, false, "Comparing incomparable type as float32: string and []uint8", "42", []byte("42"))
	testCheck(c, check.EqualsFloat32, false, "Comparing incomparable type as float32: []uint8 and string", []byte("421"), "421")
}

func (s *CheckersS) TestContains(c *check.C) {
	testCheck(c, check.Contains, false, "expected value type is not Map, Array, Slice, Chan or String", 42, 42)
	testCheck(c, check.Contains, true, "", 42, []int{42})
	testCheck(c, check.Contains, false, "expected does not contain obtained", 42, []int64{42})
	testCheck(c, check.Contains, false, "expected does not contain obtained", 42.0, []int64{11, 42, -10})
	testCheck(c, check.Contains, true, "", int64(42), []int64{12, 42, 10})
	testCheck(c, check.Contains, false, "expected does not contain obtained", int32(42), []int64{42, 10})
	testCheck(c, check.Contains, false, "expected does not contain obtained", int32(42), []int64{0, 10})

	type A struct {
		a int
		b int
	}

	testCheck(c, check.Contains, true, "", A{a: 10, b: 20}, []A{{a: 1, b: 0}, {a: 10, b: 20}})
	testCheck(c, check.Contains, true, "", &A{a: 10, b: 20}, []*A{{a: 1, b: 0}, {a: 10, b: 20}})
	testCheck(c, check.Contains, true, "", &A{a: 10, b: 20}, []*A{{a: 1, b: 0}, {a: 10, b: 20}, {a: 10, b: 20}})
	testCheck(c, check.Contains, false, "expected does not contain obtained", &A{b: 20, a: 10}, []*A{{a: 5, b: 16}, {a: 1, b: 0}, {a: 5, b: 16}})
}

func (s *CheckersS) TestNotContains(c *check.C) {
	testCheck(c, check.NotContains, true, "", 42, 42)
	testCheck(c, check.NotContains, false, "expected value contains obtained value", 42, []int{42})
	testCheck(c, check.NotContains, true, "", 42, []int64{42})
	testCheck(c, check.NotContains, true, "", 42.0, []int64{11, 42, -10})
	testCheck(c, check.NotContains, false, "expected value contains obtained value", int64(42), []int64{12, 42, 10})
	testCheck(c, check.NotContains, true, "", int32(42), []int64{42, 10})
	testCheck(c, check.NotContains, true, "", int32(42), []int64{0, 10})

	type A struct {
		a int
		b int
	}

	testCheck(c, check.NotContains, false, "expected value contains obtained value", A{a: 10, b: 20}, []A{{a: 1, b: 0}, {a: 10, b: 20}})
	testCheck(c, check.NotContains, false, "expected value contains obtained value", &A{a: 10, b: 20}, []*A{{a: 1, b: 0}, {a: 10, b: 20}})
	testCheck(c, check.NotContains, false, "expected value contains obtained value", &A{a: 10, b: 20}, []*A{{a: 1, b: 0}, {a: 10, b: 20}, {a: 10, b: 20}})
	testCheck(c, check.NotContains, true, "", &A{b: 20, a: 10}, []*A{{a: 5, b: 16}, {a: 1, b: 0}, {a: 5, b: 16}})
}

func (s *CheckersS) TestErrorIs(c *check.C) {
	e1 := errors.New("my error")
	testCheck(c, check.ErrorIs, false, "expected value is not an error", e1, 1)
	testCheck(c, check.ErrorIs, false, "obtained value is not an error", 1, e1)
	testCheck(c, check.ErrorIs, true, "", e1, e1)

	e2 := fmt.Errorf("level 1 error: %w", e1)
	testCheck(c, check.ErrorIs, true, "", e2, e1)
	e3 := fmt.Errorf("level 2 error: %w", e2)
	testCheck(c, check.ErrorIs, true, "", e3, e1)
	testCheck(c, check.ErrorIs, true, "", e3, e2)
	//

	testCheck(c, check.ErrorIs, false, "expected error doesn't contains obtained error", errors.New("my error"), e1)
	testCheck(c, check.ErrorIs, false, "expected error doesn't contains obtained error", e1, errors.New("my error"))

	testCheck(c, check.ErrorIs, false, "expected error doesn't contains obtained error", errors.New("my error"), e3)
	testCheck(c, check.ErrorIs, false, "expected error doesn't contains obtained error", e3, errors.New("my error"))
}
