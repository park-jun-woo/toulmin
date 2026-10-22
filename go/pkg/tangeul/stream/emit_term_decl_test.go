//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestEmitTermDecl — const vs struct term declarations, an `as` spec ref, and an undeclared alias error
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
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
	span := toulmin.Span{}

	cases := []struct {
		name    string
		d       ast.Definition
		wantErr bool
	}{
		{"struct def", ast.Definition{Name: "config", Kind: ast.StructDef}, false},
		{"string list def", ast.Definition{Name: "domains", Kind: ast.StringListDef, Strings: []string{"a,b", "c"}}, false},
		{"file spec def", ast.Definition{Name: "domains", Kind: ast.FileSpecDef, LoadType: "JSON list", FilePath: "./x.json"}, false},
		{"const def, no spec ref", ast.Definition{Name: "threshold", Kind: ast.ConstDef, Value: "650"}, false},
		{"const def, declared spec ref alias", ast.Definition{Name: "threshold", Kind: ast.ConstDef, Value: "650", SpecRef: &ast.Ref{Alias: "sensor", Name: "Threshold"}}, false},
		{"const def, undeclared spec ref alias", ast.Definition{Name: "threshold", Kind: ast.ConstDef, Value: "650", SpecRef: &ast.Ref{Alias: "missing", Name: "Threshold"}}, true},
	}

	check := func(t *testing.T, d ast.Definition, wantErr bool) {
		tt := &tidTable{
			byRef:       map[string]toulmin.TID{},
			aliasByName: map[string]toulmin.TID{"sensor": 2},
		}
		e := &enc{tt: tt, wire: testWire, next: 100}
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

	termRef := doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{0}}
	cdef := ast.Definition{Name: "limit", Kind: ast.ConstDef, Value: "5", SpecRef: &ast.Ref{Name: "spec"}}
	if _, err := emitTermDecl(testEnc(nodeTT(termRef, 0x10000)), termRef, cdef, toulmin.Span{}); err == nil {
		t.Fatal("node id overflow: want error")
	}
	failsAtEachID(t, 8, nodeTT(termRef, 1), func(e *enc) error {
		_, err := emitTermDecl(e, termRef, cdef, toulmin.Span{})
		return err
	})
}
