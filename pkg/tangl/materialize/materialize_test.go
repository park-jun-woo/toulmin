package materialize

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

type testLoader struct {
	calls   int
	fail    bool
	request DefinitionLoadRequest
}

func (l *testLoader) LoadDefinition(r DefinitionLoadRequest) (DefinitionSnapshot, error) {
	l.calls++
	l.request = r
	if l.fail {
		return DefinitionSnapshot{}, errors.New("broken")
	}
	return DefinitionSnapshot{Strings: []string{"a,b", "c"}}, nil
}

type testRegistry struct{ loader DefinitionLoader }

func (r testRegistry) LookupDefinitionLoader(name string) (DefinitionLoader, bool) {
	return r.loader, name == "JSON list"
}

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

func TestListsRejectsMissingRegistry(t *testing.T) {
	doc := &ast.Document{Path: "/tmp/p.md", Defs: []ast.Definition{{Name: "f", Kind: ast.FileSpecDef, LoadType: "x", FilePath: "x"}}}
	if _, err := Lists(doc, nil); err == nil {
		t.Fatal("missing registry accepted")
	}
}
