//ff:func feature=tangl type=codegen control=sequence
//ff:what TestBuildAliasImportsChecking — checking aliases join the import set; an unresolvable checking import errors
package gen

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestBuildAliasImportsChecking(t *testing.T) {
	checkingCase := ast.Case{Name: "outer", Nodes: []ast.Node{{Name: "n", Checking: &ast.Ref{Alias: "gate", Name: "can pass"}}}}

	unseen := &ast.Document{Cases: []ast.Case{checkingCase}, Rules: []ast.InlineRule{{Name: "h", Using: &ast.Ref{Alias: "nlp", Name: "f"}}}}
	specs, err := buildAliasImports(unseen)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 || specs[0].Alias != "nlp" {
		t.Fatalf("checking alias without a See is skipped, using alias kept: %+v", specs)
	}

	noModule := &ast.Document{
		Path:  filepath.Join(t.TempDir(), "doc.geul"),
		Sees:  []ast.See{{Alias: "gate", Package: "../gate/gate"}},
		Cases: []ast.Case{checkingCase},
	}
	specs, err = buildAliasImports(noModule)
	if err == nil || !strings.Contains(err.Error(), "go.mod") || specs != nil {
		t.Fatalf("specs=%+v err=%v, want a go.mod lookup error", specs, err)
	}
}
