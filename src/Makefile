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



testpackage: _test/gocheck/checkers.a

_test/gocheck/checkers.a: checkers/_test/gocheck/checkers.a
	@mkdir -p _test/gocheck
	cp checkers/_test/gocheck/checkers.a _test/gocheck

checkers/_test/gocheck/checkers.a:
	make -C checkers _test/gocheck/checkers.a
