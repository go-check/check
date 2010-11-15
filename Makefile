include $(GOROOT)/src/Make.inc

TARG=gocheck
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	gocheck.go\
	helpers.go\
	run.go\
	checkers.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w gocheck.go
	${GOFMT} -w helpers.go
	${GOFMT} -w run.go
	${GOFMT} -w checkers.go
