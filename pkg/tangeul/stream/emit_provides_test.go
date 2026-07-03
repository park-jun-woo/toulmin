//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitProvides — dispatches endpoint decl, require, run, check and rejects an unknown slotKind
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitProvides walks every path emitProvides routes: a one-element path
// (endpoint declaration), and each Path[1] slotKind (0 require, 1 run edge,
// 2 check edge, and an unrecognized slotKind, a hard error). Every switch arm
// here is a bare `return someEmitter(...)` statement, so it is fully covered
// regardless of whether the delegated emitter itself succeeds.
func TestEmitProvides(t *testing.T) {
	span := word.Span{}

	cases := []struct {
		name    string
		ref     doc.StmtRef
		elem    any
		wantErr bool
	}{
		{"endpoint decl", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0}}, ast.Endpoint{Name: "ep"}, false},
		{"require", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 0, 0}}, ast.Require{Field: "f"}, false},
		{"run edge", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 1, 0}}, "target", false},
		{"check edge", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 2, 0}}, "target", false},
		{"unknown slotKind", doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{0, 3, 0}}, "x", true},
	}

	check := func(t *testing.T, ref doc.StmtRef, elem any, wantErr bool) {
		tt := &tidTable{
			byRef:       map[string]uint16{},
			aliasByName: map[string]uint16{},
			caseByName:  map[string]uint16{"target": 9},
			nodeByCase:  map[string]uint16{},
		}
		e := &enc{tt: tt, nodeNext: 100}
		_, err := emitProvides(e, ref, elem, span)
		if wantErr && err == nil {
			t.Fatalf("emitProvides(%+v): expected error, got nil", ref)
		}
		if !wantErr && err != nil {
			t.Fatalf("emitProvides(%+v): %v", ref, err)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.ref, c.elem, c.wantErr) })
	}
}
