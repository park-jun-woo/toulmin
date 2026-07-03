//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitTermDecl — const vs struct term declarations, an `as` spec ref, and an undeclared alias error
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitTermDecl covers d.Kind != ast.ConstDef (struct terms return after the
// name label, skipping the means/sameAs edges), a const term without a
// SpecRef, a const term with a SpecRef resolved through a declared alias, and
// a SpecRef referencing an undeclared alias (a hard error from emitFuncRef).
// appendNode/appendLabel/appendEdge never fail here (NodeTerm, PropMeans and
// PropSameAs are always emitted with consistent, valid flags), so those
// error-return lines are unreachable through any legitimate call and are not
// exercised.
func TestEmitTermDecl(t *testing.T) {
	ref := doc.StmtRef{Sec: 3, Path: []int{0}}
	span := word.Span{}

	cases := []struct {
		name    string
		d       ast.Definition
		wantErr bool
	}{
		{"struct def", ast.Definition{Name: "config", Kind: ast.StructDef}, false},
		{"const def, no spec ref", ast.Definition{Name: "threshold", Kind: ast.ConstDef, Value: "650"}, false},
		{"const def, declared spec ref alias", ast.Definition{Name: "threshold", Kind: ast.ConstDef, Value: "650", SpecRef: &ast.Ref{Alias: "sensor", Name: "Threshold"}}, false},
		{"const def, undeclared spec ref alias", ast.Definition{Name: "threshold", Kind: ast.ConstDef, Value: "650", SpecRef: &ast.Ref{Alias: "missing", Name: "Threshold"}}, true},
	}

	check := func(t *testing.T, d ast.Definition, wantErr bool) {
		tt := &tidTable{
			byRef:       map[string]uint16{},
			aliasByName: map[string]uint16{"sensor": 2},
		}
		e := &enc{tt: tt, nodeNext: 100}
		got, err := emitTermDecl(e, ref, d, span)
		if wantErr && err == nil {
			t.Fatalf("emitTermDecl(%+v): expected error, got nil", d)
		}
		if wantErr {
			return
		}
		if err != nil {
			t.Fatalf("emitTermDecl(%+v): %v", d, err)
		}
		if len(got) == 0 {
			t.Errorf("emitTermDecl(%+v): expected non-empty words", d)
		}
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { check(t, c.d, c.wantErr) })
	}
}
