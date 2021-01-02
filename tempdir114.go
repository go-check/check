// +build !go1.15

package check

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// -----------------------------------------------------------------------
// Handling of temporary files and directories.

// Create a new temporary directory which is automatically removed after
// the suite finishes running.
func (c *C) MkDir() string {
	c.Helper()
	path := c.tempDir.newPath()
	if err := os.Mkdir(path, 0700); err != nil {
		c.Fatalf("Couldn't create temporary directory %s: %s", path, err.Error())
	}
	return path
}

type tempDir struct {
	sync.Mutex
	path    string
	counter int
}

func (td *tempDir) newPath() string {
	td.Lock()
	defer td.Unlock()
	if td.path == "" {
		path, err := ioutil.TempDir("", "check-")
		if err != nil {
			panic("Couldn't create temporary directory: " + err.Error())
		}
		td.path = path
	}
	result := filepath.Join(td.path, strconv.Itoa(td.counter))
	td.counter++
	return result
}

func (td *tempDir) removeAll() {
	td.Lock()
	defer td.Unlock()
	if td.path != "" {
		err := os.RemoveAll(td.path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: Error cleaning up temporaries: "+err.Error())
		}
	}
}
