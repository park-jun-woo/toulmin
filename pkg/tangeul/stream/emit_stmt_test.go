//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitStmt — dispatches every section code to its emitter and rejects an unknown section
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitStmt walks every ref.Sec case emitStmt routes on: Subject, See,
// Definitions, Rules, Cases, Provides, Internal, and an unrecognized section
// code (a hard error). Each case is a bare `return someEmitter(...)`
// statement, so every switch arm is fully covered as soon as it is reached.
func TestEmitStmt(t *testing.T) {
	span := word.Span{}

	cases := []struct {
		name    string
		ref     doc.StmtRef
		elem    any
		wantErr bool
	}{
		{"subject", doc.StmtRef{Sec: codebook.SectionSubject}, "sys", false},
		{"see", doc.StmtRef{Sec: codebook.SectionSee, Path: []int{0}}, ast.See{Alias: "a", Package: "p"}, false},
		{"definitions", doc.StmtRef{Sec: codebook.SectionDefinitions, Path: []int{0}}, ast.Definition{Name: "n", Kind: ast.StructDef}, false},
		{"rules", doc.StmtRef{Sec: codebook.SectionRules, Path: []int{0}}, ast.InlineRule{Name: "r", Cond: ast.Compare{Field: "f", Op: "equals", Value: "1"}}, false},
		{"cases", doc.StmtRef{Sec: codebook.SectionCases, Path: []int{0}}, ast.Case{Name: "c"}, false},
		{"provides", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0}}, ast.Endpoint{Name: "e"}, false},
		{"internal", doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{0}}, ast.Internal{Kind: ast.OnEvent, Event: "ev"}, false},
		{"unknown section", doc.StmtRef{Sec: codebook.SectionCode(99)}, nil, true},
	}

	check := func(t *testing.T, ref doc.StmtRef, elem any, wantErr bool) {
		e := &enc{tt: &tidTable{
			byRef:       map[string]uint16{},
			aliasByName: map[string]uint16{},
			caseByName:  map[string]uint16{},
			nodeByCase:  map[string]uint16{},
		}, nodeNext: 100}
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
