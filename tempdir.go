// +build go1.15

package check

// -----------------------------------------------------------------------
// Handling of temporary files and directories.

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (c *C) MkDir() string {
	c.Helper()
	return c.TempDir()
}

type tempDir struct{}

func (td *tempDir) removeAll() {
	// nothing
}
