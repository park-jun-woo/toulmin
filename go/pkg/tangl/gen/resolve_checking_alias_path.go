//ff:func feature=tangl type=codegen control=sequence
//ff:what resolveCheckingAliasPath — resolves a checking-alias See.Package into a sibling document path
package gen

import "path/filepath"

// resolveCheckingAliasPath resolves a checking-alias's See.Package into an
// absolute sibling document path, relative to the directory of the document
// that declared it (referencingDocPath, i.e. ast.Document.Path). This is
// gen's own copy of the same arithmetic pkg/tangeul/registry performs — kept
// here so gen never has to import upward into pkg/tangeul (§ design 4). Only
// the resolved directory matters to deriveGoImportPath, so a missing
// ".geul" extension is appended purely to yield a file-shaped path.
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
