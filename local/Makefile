include $(GOROOT)/src/Make.inc

TARG=gocheck/local
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	local.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w local.go

testpackage: _test/gocheck.a

_test/gocheck.a: ../_test/gocheck.a
	@mkdir -p _test
	cp ../_test/gocheck.a _test

../_test/gocheck.a:
	make -C .. _test/gocheck.a
