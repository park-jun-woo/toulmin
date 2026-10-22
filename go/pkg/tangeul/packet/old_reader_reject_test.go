//ff:func feature=tangl type=util control=sequence
//ff:what TestOldReaderRejectsNewRecords — a pre-Phase003 reader hits a loud decode error on each of the three new record kinds
package packet

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestOldReaderRejectsNewRecords proves §design 6's forward-incompatibility
// argument mechanically. A pre-Phase003 reader split the header as kind =
// header>>14 and count = header & 0x3FFF. Applying those legacy masks to the
// real bytes each new emitter produces, all three land in a legacy branch that
// fails loudly — never a silent misread:
//   - ExprExists (0x2001) reads as legacy kind 0 (Compare); its next word is the
//     QuantAny code, which CompareOpToAST rejects.
//   - ExprCall (0x6000) reads as legacy kind 1 (Logic); its next word is a label
//     index (0), which LogicOpToAST rejects.
//   - ExprCompareX (0xA000) reads as legacy kind 2 (Not) with count 8192, an
//     impossible child count that exhausts the record stream.
func TestOldReaderRejectsNewRecords(t *testing.T) {
	oldKind := func(h uint16) uint16 { return h >> 14 }
	oldCount := func(h uint16) int { return int(h & 0x3FFF) }

	// ExprExists → legacy Compare, bad toulmin.CompareOp.
	var l1 []string
	ex, _, _ := emitExpr(ast.Exists{Source: ast.Path{Steps: []string{"s"}}, Binder: "b", Cond: ast.Compare{Field: "x", Op: "equals", Value: "1"}}, &l1)
	if oldKind(uint16(ex[0])) != uint16(ExprCompare) {
		t.Fatalf("Exists: legacy kind = %d, want %d (Compare)", oldKind(uint16(ex[0])), ExprCompare)
	}
	if _, err := codebook.CompareOpToAST(toulmin.CompareOp(ex[1])); err == nil {
		t.Error("Exists: legacy Compare reader accepted the toulmin.QuantAny word as a toulmin.CompareOp")
	}

	// ExprCall → legacy Logic, bad toulmin.LogicOp.
	var l2 []string
	cl, _, _ := emitExpr(ast.Call{Fn: ast.Ref{Alias: "check", Name: "isSpam"}}, &l2)
	if oldKind(uint16(cl[0])) != uint16(ExprLogic) {
		t.Fatalf("Call: legacy kind = %d, want %d (Logic)", oldKind(uint16(cl[0])), ExprLogic)
	}
	if _, err := codebook.LogicOpToAST(toulmin.LogicOp(cl[1])); err == nil {
		t.Error("Call: legacy Logic reader accepted the alias label word as a toulmin.LogicOp")
	}

	// ExprCompareX → legacy Not, impossible child count.
	var l3 []string
	cx, _, _ := emitExpr(ast.Compare{Field: "ip", Op: "equals", Rhs: ast.Path{Steps: []string{"clientIp"}}}, &l3)
	if oldKind(uint16(cx[0])) != uint16(ExprNot) {
		t.Fatalf("CompareX: legacy kind = %d, want %d (Not)", oldKind(uint16(cx[0])), ExprNot)
	}
	if n := oldCount(uint16(cx[0])); n <= len(cx) {
		t.Errorf("CompareX: legacy child count = %d, want an impossible count exceeding the %d record words", n, len(cx))
	}
}
