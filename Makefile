include $(GOROOT)/src/Make.inc

TARG=launchpad.net/gocheck
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	gocheck.go\
	helpers.go\
	run.go\
	checkers.go\
	printer.go\

include $(GOROOT)/src/Make.pkg

GOFMT=gofmt -spaces=true -tabwidth=4 -tabindent=false

BADFMT=$(shell $(GOFMT) -l $(GOFILES) $(filter-out printer_test.go,$(wildcard *_test.go)))

gofmt: $(BADFMT)
	@for F in $(BADFMT); do $(GOFMT) -w $$F && echo $$F; done

ifneq ($(BADFMT),)
ifneq ($(MAKECMDGOALS),gofmt)
#$(warning WARNING: make gofmt: $(BADFMT))
endif
endif

