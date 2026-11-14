//ff:func feature=tangl type=util control=sequence
//ff:what TestInterpretStmt — decodeBody/assignPrefix/buildLabelMap errors, node- and edge-shaped statements, and the multi-primary error
package stream

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestInterpretStmt exercises every branch of interpretStmt: the decodeBody,
// assignPrefix and buildLabelMap error propagations, the len(primaries)==1
// node-shaped path (success and a buildFromNode failure), the
// len(primaries)==0 edge-shaped path (success and a buildFromEdge failure),
// and the final more-than-one-primary hard error.
func TestInterpretStmt(t *testing.T) {
	span := toulmin.Span{Off: 0, Len: 4}

	t.Run("decodeBodyError", func(t *testing.T) {
		if _, err := interpretStmt(newDec(testWire), nil, nil); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("singlePrimarySuccess", func(t *testing.T) {
		d := newDec(testWire)
		d.texts[1] = "acct"
		body := []rawPacket{
			isNodeRaw(t, toulmin.NodePacket{TID: 5, Type: toulmin.NodeNamespace, Span: span}),
			isEdgeRaw(t, toulmin.EdgePacket{Prop: toulmin.PropLabel, Src: 5, DstTID: 1, Span: span}),
		}
		blk, err := interpretStmt(d, nil, body)
		if err != nil {
			t.Fatalf("interpretStmt: %v", err)
		}
		if blk.Kind != doc.StmtBlock || blk.Ref.Sec != toulmin.SectionSubject {
			t.Errorf("got %+v, want a Subject StmtBlock", blk)
		}
		if d.ast.Subject != "acct" {
			t.Errorf("Subject = %q, want %q", d.ast.Subject, "acct")
		}
	})
	t.Run("singlePrimaryBuildFromNodeError", func(t *testing.T) {
		body := []rawPacket{isNodeRaw(t, toulmin.NodePacket{TID: 5, Type: toulmin.NodeGroup, Span: span})}
		if _, err := interpretStmt(newDec(testWire), nil, body); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("zeroPrimaryBuildFromEdgeSuccess", func(t *testing.T) {
		d := newDec(testWire)
		d.curCase = 0
		d.ast.Cases = []ast.Case{{}}
		body := []rawPacket{isEdgeRaw(t, toulmin.EdgePacket{Prop: toulmin.PropAttacks, Src: 1, DstTID: 2, Span: span})}
		blk, err := interpretStmt(d, nil, body)
		if err != nil {
			t.Fatalf("interpretStmt: %v", err)
		}
		if blk.Kind != doc.StmtBlock || blk.Ref.Sec != toulmin.SectionCases {
			t.Errorf("got %+v, want a Cases StmtBlock", blk)
		}
		if len(d.ast.Cases[0].Attacks) != 1 {
			t.Errorf("Attacks = %d, want 1", len(d.ast.Cases[0].Attacks))
		}
	})
	t.Run("zeroPrimaryBuildFromEdgeError", func(t *testing.T) {
		d := newDec(testWire)
		d.texts[2] = "label"
		// Only a label edge: primaries is empty and rel (after splitEdges) is
		// empty too, so buildFromEdge's own arity check fails.
		body := []rawPacket{isEdgeRaw(t, toulmin.EdgePacket{Prop: toulmin.PropLabel, Src: 1, DstTID: 2, Span: span})}
		if _, err := interpretStmt(d, nil, body); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("multiPrimaryError", func(t *testing.T) {
		body := []rawPacket{
			isNodeRaw(t, toulmin.NodePacket{TID: 1, Type: toulmin.NodeNamespace, Span: span}),
			isNodeRaw(t, toulmin.NodePacket{TID: 2, Type: toulmin.NodeNamespace, Span: span}),
		}
		if _, err := interpretStmt(newDec(testWire), nil, body); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("assignPrefixError", func(t *testing.T) {
		d := newDec(testWire)
		body := []rawPacket{isNodeRaw(t, toulmin.NodePacket{TID: 1, Type: toulmin.NodeRule, Span: span})}
		exprWs, err := packet.EncodeExpr(packet.ExprPacket{ID: 5, Cond: ast.Compare{Field: "f", Op: "equals", Value: "x"}}, testWire)
		if err != nil {
			t.Fatalf("EncodeExpr: %v", err)
		}
		badExpr := rawOf(t, exprWs)
		badExpr.Rec.Payload = badExpr.Rec.Payload[:len(badExpr.Rec.Payload)-1]
		if _, err := interpretStmt(d, []rawPacket{badExpr}, body); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
	t.Run("buildLabelMapError", func(t *testing.T) {
		body := []rawPacket{
			isNodeRaw(t, toulmin.NodePacket{TID: 1, Type: toulmin.NodeNamespace, Span: span}),
			isEdgeRaw(t, toulmin.EdgePacket{Prop: toulmin.PropLabel, Src: 1, DstTID: 99, Span: span}),
		}
		if _, err := interpretStmt(newDec(testWire), nil, body); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	wd := newDec(toulmin.Wire{Span: 4})
	wide, _ := toulmin.EncodeNode(toulmin.NodePacket{TID: 8, Type: toulmin.NodeNamespace, Span: toulmin.Span{Off: 1 << 32}}, wd.wire)
	if _, err := interpretStmt(wd, nil, []rawPacket{rawOf(t, wide)}); err == nil {
		t.Fatal("span past 32 bits: want error")
	}
}
