//ff:type feature=tangl type=model
//ff:what ExprPacket — a Rules condition expression carried as one flattened Payload packet
package packet

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// ExprPacket wraps a Rules inline-rule condition tree (§7.1 SECTION=4). The tree
// (ast.Compare / Logic / Not / Either) is flattened in prefix order into a
// single Payload packet (§design 4): [node count N] + a self-contained label
// table + N records in prefix order. Field/value strings are held in the packet
// as label-table entries referenced by index — the design's "Text TID" realized
// locally so the packet round-trips without external Text packets (stream dedup
// is deferred to Phase 005). EncodeExpr/DecodeExpr are an exact ast.Expr identity.
type ExprPacket struct {
	Cond ast.Expr
}
