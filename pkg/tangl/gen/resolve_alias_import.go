//ff:func feature=tangl type=codegen control=sequence
//ff:what resolveAliasImport — resolves one See alias into its importSpec by usage kind
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// resolveAliasImport resolves a single See alias into an importSpec. A
// checking alias names a sibling TANGEUL document, so its import path is
// derived from that document's location (resolveCheckingAliasPath +
// deriveGoImportPath); it returns ok=false when the alias has no See
// declaration (validate rejects that upstream, so gen just skips it). A
// "using" alias is a Go import path resolved by resolveSeePath, falling back
// to "tangl/<alias>" when no See declares it.
func resolveAliasImport(doc *ast.Document, alias string, isChecking bool) (importSpec, bool, error) {
	if isChecking {
		see, ok := findSee(doc, alias)
		if !ok {
			return importSpec{}, false, nil
		}
		sibling, err := resolveCheckingAliasPath(doc.Path, see.Package)
		if err != nil {
			return importSpec{}, false, err
		}
		path, err := deriveGoImportPath(sibling)
		if err != nil {
			return importSpec{}, false, err
		}
		return importSpec{Alias: alias, Path: path}, true, nil
	}
	path := "tangl/" + alias
	if see, ok := findSee(doc, alias); ok {
		path = see.Package
	}
	return importSpec{Alias: alias, Path: resolveSeePath(path)}, true, nil
}
