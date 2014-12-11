# Overview
The Go language provides an internal testing library, named testing, which is relatively slim due to the fact that the standard library correctness by itself is verified using it. The check package, on the other hand, expects the standard library from Go to be working correctly, and builds on it to offer a richer testing framework for libraries and applications to use.

gocheck includes features such as:

* Helpful error reporting to aid on figuring problems out (see below)
* Richer test helpers: assertions which interrupt the test immediately, deep multi-type comparisons, string matching, etc
* Suite-based grouping of tests
* Fixtures: per suite and/or per test set up and tear down
* Benchmarks integrated in the suite logic (with fixtures, etc)
* Management of temporary directories
* Panic-catching logic, with proper error reporting
* Proper counting of successes, failures, panics, missed tests, skips, etc
* Explicit test skipping
* Support for expected failures
* Verbosity flag which disables output caching (helpful to debug hanging tests, for instance)
* Multi-line string reporting for more comprehensible failures
* Inclusion of comments surrounding checks on failure reports
* Fully tested (it manages to test itself reliably)

## Example Test Suite

```go
import (
	. "github.com/masukomi/check"
	"testing"
)

// hook up gocheck into the "go test" runner
func Test(t *testing.T) { TestingT(t) }

type CoreSuite struct{} // define a test suite

var _ = Suite(&CoreSuite{})

// add a test to the suite
func (s *CoreSuite) Test_isInt(c *C) {
	c.Assert(isInt(1), IsTrue)
	c.Assert(isInt(4.0), IsFalse)
	c.Assert(isInt(int64(3)), IsTrue)
}
```


Instructions
============

Install the package with:

    go get github.com/masukomi/check

Import it with:

    import "masukomi/check"


## Assertions

* DeepEquals
	* The DeepEquals checker verifies that the obtained value is deep-equal to the expected value.  The check will work correctly even when facing slices, interfaces, and values of different types (which always fail the test).
	* Example:
	```go
	c.Assert(value, DeepEquals, 42)
	c.Assert(array, DeepEquals, []string{"hi", "there"})
	```
* DoesntPanic
	* The DoesntPanic checker verifies that calling the provided zero-argument function will not cause a panic
	* Example:
	```go
	c.Assert( func() bool { return false }, DoesntPanic)
	```
* Equals
	* The Equals checker verifies that the obtained value is equal to the expected value, according to usual Go semantics for `==`.
	* Example:
	```go
	c.Assert(value, Equals, 42)
	```
* ErrorMatches
	* The ErrorMatches checker verifies that the error value is non nil and matches the regular expression provided.
	* Example:
	```go
	c.Assert(err, ErrorMatches, "perm.*denied")
	```
* FitsTypeOf
	* The FitsTypeOf checker verifies that the obtained value is assignable to a variable with the same type as the provided sample value.
	* Example:
	```go
	c.Assert(value, FitsTypeOf, int64(0))
	c.Assert(value, FitsTypeOf, os.Error(nil))
	```
* HasLen
	* The HasLen checker verifies that the obtained value has the
provided length. In many cases this is superior to using Equals
in conjuction with the len function because in case the check
fails the value itself will be printed, instead of its length,
providing more details for figuring the problem.
	* Example:
	```go
	c.Assert(list, HasLen, 5)
	```
* Implements
	* The Implements checker verifies that the obtained value
implements the interface specified via a pointer to an interface
variable.
	* Example: 
	```go
	var e os.Error
	c.Assert(err, Implements, &e)
	```
* IsFalse
	* The IsFalse checker verifies that the obtained value is false.
	* Example:
	```go
	c.Assert(value, IsFalse)
	```
* IsNil
	* The IsNil checker tests whether the obtained value is nil.
	* Example:
	```go
	c.Assert(value, IsNil)
	```
* IsTrue
	* The IsTrue checker verifies that the obtained value is true.
	* Example:
	```go
	c.Assert(value, IsTrue)
	```
* Matches
	* The Matches checker verifies that the string provided as the obtained value (or the string resulting from obtained.String()) matches the
regular expression provided.
	* Example: 
	```go
	c.Assert(err, Matches, "perm.*denied")
	```
* NotNil
	The NotNil checker verifies that the obtained value is not nil.
	* Example:
	```go
	c.Assert(iface, NotNil)
	```
* PanicMatches
	* The PanicMatches checker verifies that calling the provided zero-argument function will cause a panic with an error value matching the regular expression provided.
	* Example:
	```go
	c.Assert(func() { f(1, 2) }, PanicMatches, `open.*: no such file or directory`)
	```
* Panics
	* The Panics checker verifies that calling the provided zero-argument function will cause a panic which is deep-equal to the provided value.
	* Example:
	```go
	c.Assert(func() { f(1, 2) }, Panics, &SomeErrorType{"BOOM"})
	```
* SliceIncludes
	* The SliceIncludes checker verifies that the provided slice includes the provided object.
	* Example:
	```go
	c.Assert(aSlice, SliceIncludes, aThing)
	```
