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
    "testing"
    "flag"
    "fmt"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface{}

// Register the given value as a test suite to be run.  Any methods starting
// with the Test prefix in the given value will be considered as a test to
// be run.
func Suite(suite interface{}) interface{} {
    allSuites = append(allSuites, suite)
    return suite
}


// -----------------------------------------------------------------------
// Public running interface.

var filterFlag = flag.String("f", "",
                             "Regular expression to select " +
                             "what to run (gocheck)")

// Run all test suites registered with the Suite() function, printing
// results to stdout, and reporting any failures back to the 'testing'
// module.
func TestingT(testingT *testing.T) {
    result := RunAll(&RunConf{Filter: *filterFlag})
    println(result.String())
    if usedDeprecatedChecks {
        println("WARNING: The Check/AssertEqual() and Check/AssertMatch() " +
                "family of functions\n         is deprecated, and will be " +
                "removed in the future.  Assert() and\n         Check() " +
                "are both more comfortable and more powerful.  See the\n" +
                "         documentation at http://labix.org/gocheck for " +
                "more details.\n")
    }
    if !result.Passed() {
        testingT.Fail()
    }
}

// Run all test suites registered with the Suite() function, using the
// given run configuration.
func RunAll(runConf *RunConf) *Result {
    result := Result{}
    for _, suite := range allSuites {
        result.Add(Run(suite, runConf))
    }
    return &result
}

// Run the given test suite using the provided run configuration.
func Run(suite interface{}, runConf *RunConf) *Result {
    runner := newSuiteRunner(suite, runConf)
    return runner.run()
}


// -----------------------------------------------------------------------
// Result methods.

func (r *Result) Add(other *Result) {
    r.Succeeded += other.Succeeded
    r.Skipped += other.Skipped
    r.Failed += other.Failed
    r.Panicked += other.Panicked
    r.FixturePanicked += other.FixturePanicked
    r.Missed += other.Missed
}

func (r *Result) Passed() bool {
    return (r.Failed == 0 && r.Panicked == 0 &&
            r.FixturePanicked == 0 && r.Missed == 0 &&
            r.RunError == nil)
}

func (r *Result) String() string {
    if r.RunError != nil {
        return "ERROR: " + r.RunError.String()
    }

    var value string
    if r.Failed == 0 && r.Panicked == 0 && r.FixturePanicked == 0 &&
       r.Missed == 0 {
        value = "OK: "
    } else {
        value = "OOPS: "
    }
    value += fmt.Sprintf("%d passed", r.Succeeded)
    if r.Skipped != 0 {
        value += fmt.Sprintf(", %d skipped", r.Skipped)
    }
    if r.Failed != 0 {
        value += fmt.Sprintf(", %d FAILED", r.Failed)
    }
    if r.Panicked != 0 {
        value += fmt.Sprintf(", %d PANICKED", r.Panicked)
    }
    if r.FixturePanicked != 0 {
        value += fmt.Sprintf(", %d FIXTURE PANICKED", r.FixturePanicked)
    }
    if r.Missed != 0 {
        value += fmt.Sprintf(", %d MISSED", r.Missed)
    }
    return value
}
