//ff:func feature=tangl type=util control=sequence
//ff:what loadDocAST — load a .tangeul file's AST via tangeul.Load
package registry

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// loadDocAST reads the .tangeul file at path (extension check, tiling gate
// via tangeul.Load) and returns its AST. It is the production loader that
// Load injects into loadRegistry; tests inject a fake in its place.
func loadDocAST(path string) (*ast.Document, error) {
	d, err := tangeul.Load(path)
	if err != nil {
		return nil, err
	}
	return d.AST, nil
}
