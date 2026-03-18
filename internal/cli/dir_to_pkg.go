//ff:func feature=cli type=util control=sequence
//ff:what dirToPkg — extracts Go package name from directory path
package cli

import (
	"path/filepath"
	"strings"
)

// dirToPkg extracts a package name from a directory path.
func dirToPkg(dir string) string {
	base := filepath.Base(dir)
	base = strings.ReplaceAll(base, "-", "")
	if base == "." || base == "" {
		return "main"
	}
	return base
}
