package check

import (
	"fmt"
	"strings"

	"github.com/kr/pretty"
)

// formatUnequal will dump the actual and expected values into a textual
// representation and return an error message containing a diff.
func formatUnequal(obtained interface{}, expected interface{}) string {
	diff := pretty.Diff(obtained, expected)
	if len(diff) == 0 {
		diff = []string{fmt.Sprintf("%p != %p", expected, obtained)}
	}

	return fmt.Sprintf(`Values are different, diff:
%s`, strings.Join(diff, "\n"))
}
