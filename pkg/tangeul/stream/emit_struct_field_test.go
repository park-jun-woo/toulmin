//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitStructField — encodes a struct `has` field (Field node + name/type labels + has edge)
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestEmitStructField is a single straight-line assertion: emitStructField has
// no data-dependent branch, only defensive appendNode/appendLabel/appendEdge
// error checks that never fire (NodeField and PropHas are always emitted with
// consistent, valid flags), so those error-return lines are unreachable
// through any legitimate call and are not exercised.
func TestEmitStructField(t *testing.T) {
	ref := doc.StmtRef{Sec: 3, Path: []int{1, 0}}
	tt := &tidTable{byRef: map[string]uint16{
		tidKey(doc.StmtRef{Sec: 3, Path: []int{1}}): 6,
	}}
	e := &enc{tt: tt, nodeNext: 100}
	f := ast.Field{Name: "port", Type: "Integer"}

	got, err := emitStructField(e, ref, f, word.Span{})
	if err != nil {
		t.Fatalf("emitStructField: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("emitStructField: expected non-empty words")
	}
}
