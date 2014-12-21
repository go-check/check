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

var testString = "Foo"

// This will run before every test
// reassigning testString to "Bar"
// regardless of what it got set to.
func (s *CoreSuite) SetUpTest(c *C) {
    testString = "Bar"
}


var _ = Suite(&CoreSuite{})

// add a test to the suite
func (s *CoreSuite) Test_isInt(c *C) {
	c.Assert(testString, Equals, "Bar")
	//       ^^^^ reset in SetUpTest
	c.Assert(isInt(1), IsTrue)
	c.Assert(isInt(4.0), IsFalse)
	c.Assert(isInt(int64(3)), IsTrue)
}
```

Note: The full list of available assertions can be found at the end of this document.


Instructions
============

Install the package with:

	go get github.com/masukomi/check

Import it with:
```go
	import (
		. "github.com/masukomi/check"
	)
```


## Using Fixtures

Fixtures are available by using one or more of the following methods in a test suite:

* `func (s *SuiteType) SetUpSuite(c *C)` \- Run once when the suite starts running.
* func (s *SuiteType) SetUpTest(c *C)` \- Run before each test or benchmark starts running.
* `func (s *SuiteType) TearDownTest(c *C)` \- Run after each test or benchmark runs.
* `func (s *SuiteType) TearDownSuite(c *C)` \- Run once after all tests or benchmarks have finished running.

Here is an example preparing some data in a temporary directory before each test runs:

```go
type Suite struct{
    dir string
}

func (s *MySuite) SetUpTest(c *C) {
    s.dir = c.MkDir()
    // Use s.dir to prepare some data.
}

func (s *MySuite) TestWithDir(c *C) {
    // Use the data in s.dir in the test.
}
```

## Adding Benchmarks

Benchmarks may be added by prefixing a method in the suite with _Benchmark_. The method will be called with the usual _*C_ argument, but unlike a normal test it is supposed to put the benchmarked logic within a loop iterating _c.N_ times.

For example:

```go
func (s *MySuite) BenchmarkLogic(c *C) {
    for i := 0; i < c.N; i++ {
        // Logic to benchmark
    }
}
```

These methods are only run when in benchmark mode, using the `-check.b` flag, and will present a result similar to the following when run:

```
PASS: myfile.go:67: MySuite.BenchmarkLogic 100000 14026 ns/op 
PASS: myfile.go:73: MySuite.BenchmarkOtherLogic 100000 21133 ns/op 
```

All the fixture methods are run as usual for a test method.

To obtain the timing for normal tests, use the `-check.v` flag instead.

## Skipping tests

Tests may be skipped with the `Skip` method within SetUpSuite, SetUpTest, or the test method itself. This allows selectively ignoring tests based on custom factors such as the architecture being run, flags provided to the test, or the availbility of resources (network, etc).

As an example, the following test suite will skip all the tests within the suite unless the _-live_ option is provided to _go test_:

```go

var live = flag.Bool("live", false, "Include live tests")

type LiveSuite struct{}

func (s *LiveSuite) SetUpSuite(c *C) {
    if !*live {
        c.Skip("-live not provided")
    }
}
```

## Running tests and output sample

Use the _go test_ tool as usual to run the tests:

```
$ go test

----------------------------------------------------------------------
FAIL: hello_test.go:16: S.TestHelloWorld

hello_test.go:17:
    c.Check(42, Equals, "42")
... obtained int = 42
... expected string = "42"

hello_test.go:18:
    c.Check(io.ErrClosedPipe, ErrorMatches, "BOOM")
... error string = "io: read/write on closed pipe"
... regex string = "BOOM"


OOPS: 0 passed, 1 FAILED
--- FAIL: hello_test.Test
FAIL
```


## Assertions and checks

gocheck uses two methods of `*C` to verify expectations on values obtained in test cases: `Assert` and `Check`. Both of these methods accept the same arguments, and the only difference between them is that when `Assert` fails, the test is interrupted immediately, while `Check` will fail the test, return `false`, and allow it to continue for further checks.

`Assert` and `Check` have the following types:

```go
func (c *C) Assert(obtained interface{}, chk Checker, ...args interface{})
func (c *C) Check(obtained interface{}, chk Checker, ...args interface{}) bool
```

They may be used as follows:

```go
func (s *S) TestSimpleChecks(c *C) {
    c.Assert(value, Equals, 42)
    c.Assert(s, Matches, "hel.*there")
    c.Assert(err, IsNil)
    c.Assert(foo, Equals, bar, Commentf("#CPUs == %d", runtime.NumCPU())
}
```

The last statement will display the provided message next to the usual debugging information, but only if the check fails.

Custom verifications may be defined by implementing the `Checker` interface. There are several standard checkers available. See the documtation for details and examples:

## Selecting which tests to run

gocheck can filter tests out based on the test name, the suite name, or both. To run tests selectively, provide the command line option `-check.f` when running `go test`. Note that this option is specific to `gocheck`, and won't affect `go test` itself.

Some examples:

```shell
$ go test -check.f MyTestSuite
$ go test -check.f "Test.*Works"
$ go test -check.f "MyTestSuite.Test.*Works```


## Verbose modes

gocheck offers two levels of verbosity through the `-check.v` and `-check.vv` flags. In the first mode, passing tests will also be reported. The second mode will disable log caching entirely and will stream starting and ending suite calls and everything logged in between straight to the output. This is useful to debug hanging tests, for instance.




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
-----

## License
Distributed under the [Simplified BSD License](http://en.wikipedia.org/wiki/BSD_licenses#2-clause_license_.28.22Simplified_BSD_License.22_or_.22FreeBSD_License.22.29)


-----
## Credit where credit is due:

This repo is an extension of the original [gocheck](http://labix.org/gocheck)
  ([GitHub repo](https://github.com/go-check/check#readme)) by [Gustavo Niemeyer](http://labix.org/)  
