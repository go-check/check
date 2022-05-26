package formatters

import (
	"fmt"
)

func DefaultOutput(prefix, label, funcPath, funcName, suffix string) string {
	return fmt.Sprintf("%s%s: %s: %s%s", prefix, label, funcPath, funcName, suffix)
}
