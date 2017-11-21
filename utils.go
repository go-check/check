package check

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
)

// formatUnequal will dump the actual and expected values into a textual
// representation and return an error message containing a diff.
func formatUnequal(actual interface{}, expected interface{}) string {
	a := spew.Sdump(actual)
	e := spew.Sdump(expected)

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(e),
		B:        difflib.SplitLines(a),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})
	return "Values are different, diff:\n" + diff
}
