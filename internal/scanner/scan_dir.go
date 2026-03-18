//ff:func feature=scanner type=scanner control=iteration dimension=1
//ff:what ScanDir — collects Go source file paths from a directory
package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanDir returns paths of Go source files in dir,
// excluding _test.go and _gen.go files.
func ScanDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		if strings.HasSuffix(name, "_gen.go") {
			continue
		}
		paths = append(paths, filepath.Join(dir, name))
	}
	return paths, nil
}
