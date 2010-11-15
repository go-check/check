include $(GOROOT)/src/Make.inc

TARG=gocheck
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	gocheck.go\
	helpers.go\
	run.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w gocheck.go
	${GOFMT} -w helpers.go
	${GOFMT} -w run.go



testpackage: _test/gocheck/local.a

_test/gocheck/local.a: local/_test/gocheck/local.a
	@mkdir -p _test/gocheck
	cp local/_test/gocheck/local.a _test/gocheck

local/_test/gocheck/local.a:
	make -C local _test/gocheck/local.a
