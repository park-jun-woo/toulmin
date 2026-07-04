//ff:func feature=tangl type=validator control=sequence
//ff:what checkingEdgeTarget — resolves a node's checking target into a cross-document graph key
package registry

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkingEdgeTarget returns the "<path>#<case>" graph key a node's checking
// clause points at, or "" when the node has no checking clause (or an alias
// with no See declaration, which single-document validate rejects upstream).
// A same-document target keeps docPath; an alias-qualified target resolves to
// its sibling document's canonical path.
func checkingEdgeTarget(d *ast.Document, docPath string, n ast.Node) (string, error) {
	if n.Checking == nil {
		return "", nil
	}
	if n.Checking.Alias == "" {
		return docPath + "#" + n.Checking.Name, nil
	}
	see, ok := findSee(d, n.Checking.Alias)
	if !ok {
		return "", nil
	}
	child, err := resolveCheckingAliasPath(d.Path, see.Package)
	if err != nil {
		return "", err
	}
	return child + "#" + n.Checking.Name, nil
}
