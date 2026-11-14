//ff:func feature=tangl type=engine control=sequence
//ff:what TestListsMaterializesOnceRelativeToDocument — file definitions load once with a document-relative path
package materialize

import (
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

func TestListsMaterializesOnceRelativeToDocument(t *testing.T) {
	l := &testLoader{}
	doc := &ast.Document{Path: filepath.Join(t.TempDir(), "policy.md"), Defs: []ast.Definition{
		{Name: "inline", Kind: ast.StringListDef, Strings: []string{"x"}, Line: 1},
		{Name: "file", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "data.json", Line: 2},
	}}
	got, err := Lists(doc, testRegistry{l})
	if err != nil {
		t.Fatal(err)
	}
	if l.calls != 1 || l.request.ResolvedPath != filepath.Join(filepath.Dir(doc.Path), "data.json") {
		t.Fatalf("calls/request = %d %#v", l.calls, l.request)
	}
	if got["file"].SpecName() != "file" || len(got["inline"].StringList()) != 1 {
		t.Fatalf("got = %#v", got)
	}
}
