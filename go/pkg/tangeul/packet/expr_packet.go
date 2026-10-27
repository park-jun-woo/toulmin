//ff:type feature=tangl type=model
//ff:what ExprPacket — a Rules condition expression carried as one PREDICATE-DECLARE record
package packet

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// ExprPacket wraps a Rules inline-rule condition tree, written as a
// PREDICATE-DECLARE record (GEUL t002 TOULMIN EXPR): its Edge Id, [node count N]
// + a self-contained label table of UTF8-TEXT Values + N records in prefix order.
// Field/value/operand strings are held in the packet as label-table entries
// referenced by index — the design's "Text TID" realized locally so the packet
// round-trips without external Text packets (stream dedup is deferred to
// Phase 005). EncodeExpr/DecodeExpr are an exact ast.Expr identity.
//
// Record kinds carried by a 3-bit header (packExprHeader, §design 1): the legacy
// ExprCompare / Logic / Not / Either, plus the Phase003 additions ExprExists (an
// "any of … satisfies" quantifier: [QuantOp][source operand][binder] + one Cond
// child), ExprCall (a predicate host call: [alias][name]) and ExprCompareX (a
// comparison whose Lhs and/or Rhs is an operand: [op][field][lhs present +
// operand?][value | ValueOperand + operand]). Operands (ast.Path / SpecRef /
// Call) ride as self-contained sub-records (§design 2). Emission rule preserving
// backward compatibility: an ast.Compare with Lhs == nil && Rhs == nil emits the
// legacy ExprCompare record byte-for-byte as before; ExprCompareX is used only
// when an operand is present, so documents without the new syntax are unchanged.
type ExprPacket struct {
	ID   toulmin.TID
	Cond ast.Expr
}
