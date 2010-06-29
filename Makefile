include $(GOROOT)/src/Make.$(GOARCH)

TARG=gocheck
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	gocheck.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w gocheck.go
