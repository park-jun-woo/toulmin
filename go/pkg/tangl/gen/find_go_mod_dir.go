//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what findGoModDir — walks up from dir to the nearest directory containing a go.mod
package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// findGoModDir walks up from dir until it finds a directory holding a go.mod,
// returning that directory (the module root). It stops at the filesystem root
// with an error when no enclosing module exists.
func findGoModDir(dir string) (string, error) {
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("gen: no go.mod found above %q", dir)
		}
		dir = parent
	}
}
