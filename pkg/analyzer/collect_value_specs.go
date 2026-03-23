//ff:func feature=analyzer type=analyzer control=iteration dimension=2
//ff:what collectValueSpecs — collects ValueSpec nodes from GenDecl in AST file
package analyzer

import "go/ast"

// collectValueSpecs extracts all ValueSpec nodes from top-level GenDecls.
func collectValueSpecs(f *ast.File) []*ast.ValueSpec {
	var specs []*ast.ValueSpec
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			specs = append(specs, vs)
		}
	}
	return specs
}
