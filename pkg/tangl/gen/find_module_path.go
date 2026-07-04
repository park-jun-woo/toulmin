//ff:func feature=tangl type=codegen control=sequence
//ff:what findModulePath — reads the module path declared by the nearest enclosing go.mod
package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// findModulePath locates the nearest go.mod at or above dir (findGoModDir) and
// returns that module root directory and the module path declared on its
// "module" line (moduleLine). It mirrors how the go tool resolves the
// enclosing module, so a generated package's import path can be derived from
// where its source .tangeul lives — even in a module other than this repo's.
func findModulePath(dir string) (root string, module string, err error) {
	root, err = findGoModDir(dir)
	if err != nil {
		return "", "", err
	}
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", "", err
	}
	module, ok := moduleLine(b)
	if !ok {
		return "", "", fmt.Errorf("gen: %s has no module line", filepath.Join(root, "go.mod"))
	}
	return root, module, nil
}
