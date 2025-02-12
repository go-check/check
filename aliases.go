package check

// File contents aliases for Checker over Assert

const aliasSkippedFrame = 1

func (c *C) runObtainedExpectedAlias(funcName string, obtained, expected any, checker Checker, args any) bool {
	comf := commentArgs(args)

	if comf != nil {
		if !c.internalCheck(aliasSkippedFrame, funcName, obtained, checker, expected, commentArgs(args)) {
			c.stopNow()
			return false
		}
	} else if !c.internalCheck(aliasSkippedFrame, funcName, obtained, checker, expected) {
		c.stopNow()
		return false
	}

	return true
}

func (c *C) runObtainedAlias(funcName string, obtained any, checker Checker, args any) bool {
	comf := commentArgs(args)

	if comf != nil {
		if !c.internalCheck(aliasSkippedFrame, funcName, obtained, checker, commentArgs(args)) {
			c.stopNow()
			return false
		}
	} else if !c.internalCheck(aliasSkippedFrame, funcName, obtained, checker) {
		c.stopNow()
		return false
	}

	return true
}

func (c *C) Contains(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Contains", obtained, expected, Contains, args)
}

func (c *C) NotContains(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotContains", obtained, expected, NotContains, args)
}

func (c *C) DeepEquals(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("DeepEquals", obtained, expected, DeepEquals, args)
}

func (c *C) NotDeepEquals(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotDeepEquals", obtained, expected, Not(DeepEquals), args)
}

func (c *C) Equals(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Equals", obtained, expected, Equals, args)
}

func (c *C) NotEquals(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotEquals", obtained, expected, NotEquals, args)
}

func (c *C) EqualsFloat32(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("EqualsFloat32", obtained, expected, EqualsFloat32, args)
}

func (c *C) NotEqualsFloat32(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotEqualsFloat32", obtained, expected, Not(EqualsFloat32), args)

}

func (c *C) EqualsMore(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("EqualsMore", obtained, expected, EqualsMore, args)
}

func (c *C) NotEqualsMore(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotEqualsMore", obtained, expected, Not(EqualsMore), args)
}

func (c *C) ErrorMatches(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("ErrorMatches", obtained, expected, ErrorMatches, args)
}

func (c *C) FitsTypeOf(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("FitsTypeOf", obtained, expected, FitsTypeOf, args)
}

func (c *C) Len(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Len", obtained, expected, HasLen, args)
}

func (c *C) HasLen(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("HasLen", obtained, expected, HasLen, args)
}

func (c *C) NotHasLen(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotHasLen", obtained, expected, Not(HasLen), args)
}

func (c *C) HasLenLessThan(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("HasLenLessThan", obtained, expected, HasLenLessThan, args)
}

func (c *C) HasLenMoreThan(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("HasLenMoreThan", obtained, expected, HasLenMoreThan, args)
}

func (c *C) Implements(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Implements", obtained, expected, Implements, args)
}

func (c *C) ErrorIs(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("ErrorIs", obtained, expected, ErrorIs, args)
}

func (c *C) False(obtained any, args ...any) bool {
	return c.runObtainedAlias("False", obtained, IsFalse, args)
}

func (c *C) IsFalse(obtained any, args ...any) bool {
	return c.runObtainedAlias("IsFalse", obtained, IsFalse, args)
}

func (c *C) True(obtained any, args ...any) bool {
	return c.runObtainedAlias("True", obtained, IsTrue, args)
}

func (c *C) IsTrue(obtained any, args ...any) bool {
	return c.runObtainedAlias("IsTrue", obtained, IsTrue, args)
}

func (c *C) NotNil(obtained any, args ...any) bool {
	return c.runObtainedAlias("NotNil", obtained, NotNil, args)
}

func (c *C) Equal(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Equal", obtained, expected, Equals, args)
}

func (c *C) NotEqual(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotEqual", obtained, expected, NotEquals, args)
}

func (c *C) Nil(obtained any, args ...any) bool {
	return c.runObtainedAlias("Nil", obtained, IsNil, args)
}

func (c *C) IsNil(obtained any, args ...any) bool {
	return c.runObtainedAlias("IsNil", obtained, IsNil, args)
}

func (c *C) Empty(obtained any, args ...any) bool {
	return c.runObtainedAlias("Empty", obtained, IsEmpty, args)
}

func (c *C) IsEmpty(obtained any, args ...any) bool {
	return c.runObtainedAlias("IsEmpty", obtained, IsEmpty, args)
}

func (c *C) NotEmpty(obtained any, args ...any) bool {
	return c.runObtainedAlias("NotEmpty", obtained, NotEmpty, args)
}

func (c *C) LessOrEqual(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("LessOrEqual", obtained, expected, LessOrEqualThan, args)
}

func (c *C) LessOrEqualThan(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("LessOrEqualThan", obtained, expected, LessOrEqualThan, args)
}

func (c *C) Less(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Less", obtained, expected, LessThan, args)
}

func (c *C) LessThan(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("LessThan", obtained, expected, LessThan, args)
}

func (c *C) Matches(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Matches", obtained, expected, Matches, args)
}

func (c *C) NotMatches(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("NotMatches", obtained, expected, Not(Matches), args)
}

func (c *C) GreaterOrEqual(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("GreaterOrEqual", obtained, expected, MoreOrEqualThan, args)
}

func (c *C) MoreOrEqualThan(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("MoreOrEqualThan", obtained, expected, MoreOrEqualThan, args)
}

func (c *C) Greater(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Greater", obtained, expected, MoreThan, args)
}

func (c *C) MoreThan(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("MoreThan", obtained, expected, MoreThan, args)
}

func (c *C) PanicMatches(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("PanicMatches", obtained, expected, PanicMatches, args)
}

func (c *C) Panics(obtained, expected any, args ...any) bool {
	return c.runObtainedExpectedAlias("Panics", obtained, expected, Panics, args)
}
