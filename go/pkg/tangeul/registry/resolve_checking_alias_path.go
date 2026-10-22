//ff:func feature=tangl type=util control=sequence
//ff:what resolveCheckingAliasPath — resolves a checking-alias's See.Package into an absolute .geul path
package registry

import "path/filepath"

// resolveCheckingAliasPath resolves a checking-alias's See.Package into an
// absolute .geul file path, relative to the directory of the document
// that declared it (referencingDocPath, i.e. ast.Document.Path). Unlike
// gen's resolveSeePath (Go-import resolution for `using`), Package here is a
// filesystem path: a missing ".geul" extension is appended, and the
// result is resolved to a cleaned absolute path.
func resolveCheckingAliasPath(referencingDocPath, pkg string) (string, error) {
	p := pkg
	if filepath.Ext(p) != ".geul" {
		p += ".geul"
	}
	if !filepath.IsAbs(p) {
		p = filepath.Join(filepath.Dir(referencingDocPath), p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}
