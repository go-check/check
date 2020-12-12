// +build go1.15

package check

type tempDir struct{}

func (td *tempDir) removeAll() {
	// nothing
}
