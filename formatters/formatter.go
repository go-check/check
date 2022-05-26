package formatters

import (
	"strings"
	"time"
)

type Formatter string

const (
	DefaultFormatter  = "default"
	JsonFormatter     = "json"
	TeamcityFormatter = "teamcity"
)

func F(name *string) Formatter {
	switch strings.ToLower(*name) {
	case "json":
		return JsonFormatter
	case "teamcity":
		return TeamcityFormatter
	}

	return DefaultFormatter
}

type Data struct {
	StartTime time.Time
	Duration  time.Duration

	TestName string
	StdOut   string
	FuncPath string
	Package  string
	FuncName string
	Prefix   string
	Label    string
	Suffix   string

	Formatter    Formatter
	FormatPrefix string
}

func Render(d Data) string {
	out := ""

	switch d.Formatter {
	case TeamcityFormatter:
		out = DefaultOutput(d.Prefix, d.Label, d.FuncPath, d.FuncName, d.Suffix) +
			TeamcityOutput(d.Label, d.TestName, d.StdOut, d.StartTime, d.Duration,
				d.FuncPath, d.FuncName, d.FormatPrefix, d.Suffix) + "\n"
	case JsonFormatter:
		out = JsonOutput(d.Label, d.TestName, d.StdOut, d.StartTime, d.Duration,
			d.FuncPath) + "\n"
	default:
		out = DefaultOutput(d.Prefix, d.Label, d.FuncPath, d.FuncName, d.Suffix)
	}

	return out
}
