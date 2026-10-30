//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildImportsAliasError — a failing alias import resolution aborts buildImports with no specs
package gen

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestBuildImportsAliasError(t *testing.T) {
	doc := &ast.Document{
		Path:  filepath.Join(t.TempDir(), "doc.geul"),
		Sees:  []ast.See{{Alias: "gate", Package: "gate"}},
		Cases: []ast.Case{{Name: "c", Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "gate", Name: "x"}}}}},
	}
	specs, err := buildImports(&genContext{Doc: doc})
	if err == nil || specs != nil {
		t.Fatalf("specs=%+v err=%v, want the alias import error", specs, err)
	}
}
