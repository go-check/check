package gocheck

import (
    "testing"
    "os"
    "io"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface{}

func Suite(suite interface{}) interface{} {
    return Suites(suite)[0]
}

func Suites(suites ...interface{}) []interface{} {
    lenAllSuites := len(allSuites)
    lenSuites := len(suites)
    if lenAllSuites + lenSuites > cap(allSuites) {
        newAllSuites := make([]interface{}, (lenAllSuites+lenSuites)*2)
        copy(newAllSuites, allSuites)
        allSuites = newAllSuites
    }
    allSuites = allSuites[0:lenAllSuites+lenSuites]
    for i, suite := range suites {
        allSuites[lenAllSuites+i] = suite
    }
    return suites
}


// -----------------------------------------------------------------------
// Test running logic.

func TestingT(testingT *testing.T) {
    RunAll()
}

func RunAll() {
    for _, suite := range allSuites {
        Run(suite)
    }
}

func Run(suite interface{}) {
    RunWithWriter(suite, os.Stdout)
}

func RunWithWriter(suite interface{}, writer io.Writer) {
    runner := newSuiteRunner(suite, writer)
    runner.run()
}
