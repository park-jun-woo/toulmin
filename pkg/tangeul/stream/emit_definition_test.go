//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitDefinition — a two-element path routes to a struct field, any other path to the term decl
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitDefinition checks both emitDefinition branches: len(Path)==2 (a
// struct `has` field) dispatches to emitStructField, and any other Path length
// (the term declaration itself, Path==[i]) dispatches to emitTermDecl.
func TestEmitDefinition(t *testing.T) {
	e := &enc{tt: &tidTable{byRef: map[string]uint16{}}}
	span := word.Span{TID: 0, Off: 0, Len: 1}

	t.Run("struct field", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 3, Path: []int{1, 0}}
		got, err := emitDefinition(e, ref, ast.Field{Name: "host", Type: "Text"}, span)
		if err != nil {
			t.Fatalf("emitDefinition field: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})

	t.Run("term decl", func(t *testing.T) {
		ref := doc.StmtRef{Sec: 3, Path: []int{0}}
		got, err := emitDefinition(e, ref, ast.Definition{Name: "x", Kind: ast.ConstDef, Value: "1"}, span)
		if err != nil {
			t.Fatalf("emitDefinition decl: %v", err)
		}
		if len(got) == 0 {
			t.Error("expected non-empty words")
		}
	})
}
