include $(GOROOT)/src/Make.inc

TARG=gocheck
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	gocheck.go\
	gocheckrun.go\
	gochecktest.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w gocheck.go
	${GOFMT} -w gocheckrun.go
	${GOFMT} -w gochecktest.go
