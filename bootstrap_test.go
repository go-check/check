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

// These initial tests are for bootstrapping.  They verify that we can
// basically use the testing infrastructure itself to check if the test
// system is working.
// 
// These tests use will break down the test runner badly in case of
// errors because if they simply fail, we can't be sure the developer
// will ever see anything (because failing means the failing system
// somehow isn't working! :-)
//
// Do not assume *any* internal functionality works as expected besides
// what's actually tested here.

package gocheck_test


import (
    "gocheck"
    "strings"
    "fmt"
)


type BootstrapS struct{}

var boostrapS = gocheck.Suite(&BootstrapS{})

func (s *BootstrapS) TestCountSuite(c *gocheck.C) {
    suitesRun += 1
}

func (s *BootstrapS) TestFailedAndFail(c *gocheck.C) {
    if c.Failed() {
        critical("c.Failed() must be false first!")
    }
    c.Fail()
    if !c.Failed() {
        critical("c.Fail() didn't put the test in a failed state!")
    }
    c.Succeed()
}

func (s *BootstrapS) TestFailedAndSucceed(c *gocheck.C) {
    c.Fail()
    c.Succeed()
    if c.Failed() {
        critical("c.Succeed() didn't put the test back in a non-failed state")
    }
}

func (s *BootstrapS) TestLogAndGetTestLog(c *gocheck.C) {
    c.Log("Hello there!")
    log := c.GetTestLog()
    if log != "Hello there!\n" {
        critical(fmt.Sprintf("Log() or GetTestLog() is not working! Got: %#v", log))
    }
}

func (s *BootstrapS) TestLogfAndGetTestLog(c *gocheck.C) {
    c.Logf("Hello %v", "there!")
    log := c.GetTestLog()
    if log != "Hello there!\n" {
        critical(fmt.Sprintf("Logf() or GetTestLog() is not working! Got: %#v", log))
    }
}

func (s *BootstrapS) TestRunShowsErrors(c *gocheck.C) {
    output := String{}
    gocheck.Run(&FailHelper{}, &gocheck.RunConf{Output: &output})
    if strings.Index(output.value, "Expected failure!") == -1 {
        critical(fmt.Sprintf("RunWithWriter() output did not contain the " +
                             "expected failure! Got: %#v", output.value))
    }
}

func (s *BootstrapS) TestRunDoesntShowSuccesses(c *gocheck.C) {
    output := String{}
    gocheck.Run(&SuccessHelper{}, &gocheck.RunConf{Output: &output})
    if strings.Index(output.value, "Expected success!") != -1 {
        critical(fmt.Sprintf("RunWithWriter() output contained a successful " +
                             "test! Got: %#v", output.value))
    }
}
