//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitStmt — dispatches every section code to its emitter and rejects an unknown section
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestEmitStmt walks every ref.Sec case emitStmt routes on: Subject, See,
// Definitions, Rules, Cases, Provides, Internal, and an unrecognized section
// code (a hard error). Each case is a bare `return someEmitter(...)`
// statement, so every switch arm is fully covered as soon as it is reached.
func TestEmitStmt(t *testing.T) {
	span := toulmin.Span{}

	cases := []struct {
		name    string
		ref     doc.StmtRef
		elem    any
		wantErr bool
	}{
		{"subject", doc.StmtRef{Sec: toulmin.SectionSubject}, "sys", false},
		{"see", doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{0}}, ast.See{Alias: "a", Package: "p"}, false},
		{"definitions", doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{0}}, ast.Definition{Name: "n", Kind: ast.StructDef}, false},
		{"rules", doc.StmtRef{Sec: toulmin.SectionRules, Path: []int{0}}, ast.InlineRule{Name: "r", Cond: ast.Compare{Field: "f", Op: "equals", Value: "1"}}, false},
		{"cases", doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}}, ast.Case{Name: "c"}, false},
		{"provides", doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{0}}, ast.Endpoint{Name: "e"}, false},
		{"internal", doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{0}}, ast.Internal{Kind: ast.OnEvent, Event: "ev"}, false},
		{"unknown section", doc.StmtRef{Sec: toulmin.SectionCode(99)}, nil, true},
	}

	check := func(t *testing.T, ref doc.StmtRef, elem any, wantErr bool) {
		e := &enc{tt: &tidTable{
			byRef:       map[string]toulmin.TID{},
			aliasByName: map[string]toulmin.TID{},
			caseByName:  map[string]toulmin.TID{},
			nodeByCase:  map[string]toulmin.TID{},
		}, wire: testWire, next: 100}
		_, err := emitStmt(e, ref, elem, span)
		if wantErr && err == nil {
			t.Fatalf("emitStmt(%+v): expected error, got nil", ref)
		}
		if !wantErr && err != nil {
			t.Fatalf("emitStmt(%+v): %v", ref, err)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.ref, c.elem, c.wantErr) })
	}
}
