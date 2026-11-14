//ff:func feature=tangl type=codegen control=sequence
//ff:what deriveGoImportPath — derives a resolved document's generated-package Go import path from its location
package gen

import (
	"path"
	"path/filepath"
)

// deriveGoImportPath finds the nearest go.mod above tangeulPath's directory
// and joins the module's declared path with the directory's path relative to
// the module root — mirroring how tangl gen's own output is expected to live
// alongside its source .geul (the generated .go file sits in the same
// directory, same Go package). This generalizes resolveSeePath's hard-coded
// "tangl/" convention to any module via go.mod lookup (§ design 4).
func deriveGoImportPath(tangeulPath string) (string, error) {
	abs, err := filepath.Abs(tangeulPath)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(abs)
	root, module, err := findModulePath(dir)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return module, nil
	}
	return path.Join(module, filepath.ToSlash(rel)), nil
}
