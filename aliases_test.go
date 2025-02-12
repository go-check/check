package check_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/iostrovok/check"
)

type AliasSuite struct{}

var _ = check.Suite(&AliasSuite{})

func (s *AliasSuite) TestCountSuite(c *check.C) {
	suitesRun += 1
}

func TestAlias(t *testing.T) { check.TestingT(t) }

func (s *AliasSuite) TestContainsAlias(c *check.C) {
	c.Contains(42, []int{42})
	c.Contains(42, []int{42}, "test syntax")
}

func (s *AliasSuite) TestNotContainsAlias(c *check.C) {
	c.NotContains(43, []int{42})
	c.NotContains(43, []int{42}, "test syntax")
}

func (s *AliasSuite) TestDeepEqualsAlias(c *check.C) {
	c.DeepEquals([]int{42}, []int{42})
	c.DeepEquals([]int{42}, []int{42}, "test syntax")
}

func (s *AliasSuite) TestNotDeepEqualsAlias(c *check.C) {
	c.NotDeepEquals([]int{43}, []int{42})
	c.NotDeepEquals([]int{42}, []int{42, 1}, "test syntax")
}

func (s *AliasSuite) TestEqualsAlias(c *check.C) {
	c.Equals(42, 42)
	c.Equals(42, 42, "test syntax")

	c.Equal(42, 42)
	c.Equal(42, 42, "test syntax")
}

func (s *AliasSuite) TestNotEqualsAlias(c *check.C) {
	c.NotEquals(41, 43)
	c.NotEquals(41, 43, "test syntax")

	c.NotEqual(41, 43)
	c.NotEqual(41, 43, "test syntax")
}

func (s *AliasSuite) TestEqualsFloat32Alias(c *check.C) {
	c.EqualsFloat32(0.1, 0.1)
	c.EqualsFloat32(0.1, 0.1, "test syntax")
}

func (s *AliasSuite) TestNotEqualsFloat32Alias(c *check.C) {
	c.NotEqualsFloat32(0.2, 0.1)
	c.NotEqualsFloat32(0.1, 0.3, "test syntax")
}

func (s *AliasSuite) TestEqualsMoreAlias(c *check.C) {
	c.EqualsMore(0.1, 0.1)
	c.EqualsMore(0.1, 0.1, "test syntax")
}

func (s *AliasSuite) TestNotEqualsMoreAlias(c *check.C) {
	c.NotEqualsMore(0.2, 0.1)
	c.NotEqualsMore(0.3, 0.1, "test syntax")
}

// 	testCheck(c, check.ErrorMatches, true, "", errors.New("some error"), "some error")
//	testCheck(c, check.ErrorMatches, true, "", errors.New("some error"), "so.*or")

func (s *AliasSuite) TestErrorMatchesAlias(c *check.C) {
	c.ErrorMatches(errors.New("some error"), "some error")
	c.ErrorMatches(errors.New("some error"), "some error", "test syntax")
}

func (s *AliasSuite) TestFitsTypeOfAlias(c *check.C) {
	c.FitsTypeOf(0.1, 0.1)
	c.FitsTypeOf(0.1, 0.1, "test syntax")
}

func (s *AliasSuite) TestHasLenAlias(c *check.C) {
	c.HasLen([]int{42}, 1)
	c.HasLen([]int{42, 3, 2, 3}, 4, "test syntax")

	c.Len([]int{42}, 1)
	c.Len([]int{42, 3, 2, 3}, 4, "test syntax")
}

func (s *AliasSuite) TestNotHasLenAlias(c *check.C) {
	c.NotHasLen([]int{42}, 21)
	c.NotHasLen([]int{42, 3, 2, 3}, 5, "test syntax")
}

func (s *AliasSuite) TestHasLenLessThanAlias(c *check.C) {
	c.HasLenLessThan([]int{42}, 21)
	c.HasLenLessThan([]int{42, 3, 2, 3}, 5, "test syntax")
}

func (s *AliasSuite) TestHasLenMoreThanAlias(c *check.C) {
	c.HasLenMoreThan([]int{}, -1)
	c.HasLenMoreThan([]int{42}, 0)
	c.HasLenMoreThan([]int{42, 3, 2, 3}, 2, "test syntax")
}

func (s *AliasSuite) TestImplementsAlias(c *check.C) {
	var e error

	c.Implements(errors.New(""), &e)
	c.Implements(errors.New(""), &e, "test syntax")
}

func (s *AliasSuite) TestErrorIsAlias(c *check.C) {
	e1 := errors.New("my error")
	e2 := fmt.Errorf("level 1 error: %w", e1)

	c.ErrorIs(e2, e1)
	c.ErrorIs(e2, e1, "test syntax")
}

func (s *AliasSuite) TestIsFalseAlias(c *check.C) {
	c.IsFalse(false)
	c.IsFalse(false, "test syntax")

	c.False(false)
	c.False(false, "test syntax")
}

func (s *AliasSuite) TestPanicsAlias(c *check.C) {
	c.Panics(func() { panic("BOOM") }, "BOOM")
	c.Panics(func() { panic("BOOM") }, "BOOM", "test syntax")
}

func (s *AliasSuite) TestPanicMatchesAlias(c *check.C) {
	f := func() { panic("BOOM") }
	c.PanicMatches(f, "BOOM")
	c.PanicMatches(f, "BOOM", "test syntax")
}

func (s *AliasSuite) TestNotNilAlias(c *check.C) {
	c.NotNil(false, nil)
	c.NotNil(true, nil, "test syntax")
}

func (s *AliasSuite) TestIsNilAlias(c *check.C) {
	c.IsNil(nil)
	c.IsNil(nil, "test syntax")

	c.Nil(nil)
	c.Nil(nil, "test syntax")
}

func (s *AliasSuite) TestIsTrueAlias(c *check.C) {
	c.IsTrue(true)
	c.IsTrue(true, "test syntax")

	c.True(true)
	c.True(true, "test syntax")
}

func (s *AliasSuite) TestMoreOrEqualThanAlias(c *check.C) {
	c.MoreOrEqualThan(0.1, 0.1)
	c.MoreOrEqualThan(0.1, 0.1, "test syntax")

	c.GreaterOrEqual(0.1, 0.1)
	c.GreaterOrEqual(0.1, 0.1, "test syntax")
}

func (s *AliasSuite) TestMoreThanAlias(c *check.C) {
	c.MoreThan(0.4, 0.1)
	c.MoreThan(10, 8, "test syntax")

	c.Greater(0.4, 0.1)
	c.Greater(10, 8, "test syntax")
}

func (s *AliasSuite) TestLessThanAlias(c *check.C) {
	c.LessThan(0.1, 0.4)
	c.LessThan(-10, 1024, "test syntax")

	c.Less(0.1, 0.4)
	c.Less(-10, 1024, "test syntax")
}

func (s *AliasSuite) TestLessOrEqualThanAlias(c *check.C) {
	c.LessOrEqualThan(0.1, 0.1)
	c.LessOrEqualThan(0.1, 0.1, "test syntax")

	c.LessOrEqual(0.1, 0.1)
	c.LessOrEqual(0.1, 0.1, "test syntax")
}

func (s *AliasSuite) TestNotMatchesAlias(c *check.C) {
	c.NotMatches(reflect.ValueOf("!abc"), "a.c")
	c.NotMatches(reflect.ValueOf("!abc"), "a.c", "test syntax")
}

func (s *AliasSuite) TestMatchesAlias(c *check.C) {
	c.Matches(reflect.ValueOf("abc"), "a.c")
	c.Matches(reflect.ValueOf("abc"), "a.c", "test syntax")
}

func (s *AliasSuite) TestIsEmptyAlias(c *check.C) {
	c.IsEmpty(map[int]map[string]bool{})
	c.IsEmpty("", "test syntax")

	c.Empty(map[int]map[string]bool{})
	c.Empty("", "test syntax")
}
