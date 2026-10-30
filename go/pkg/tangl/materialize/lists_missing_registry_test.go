//ff:func feature=tangl type=engine control=sequence
//ff:what TestListsRejectsMissingRegistry — a file definition without a loader registry is an error
package materialize

import (
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestListsRejectsMissingRegistry(t *testing.T) {
	doc := &ast.Document{Path: "/tmp/p.md", Defs: []ast.Definition{{Name: "f", Kind: ast.FileSpecDef, LoadType: "x", FilePath: "x"}}}
	if _, err := Lists(doc, nil); err == nil {
		t.Fatal("missing registry accepted")
	}
}
