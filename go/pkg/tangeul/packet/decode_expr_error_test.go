//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestDecodeExprErrors — DecodeExpr rejects wrong records, truncated fields, bad labels/records, leftovers and a count mismatch
package packet

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestDecodeExprErrors exercises DecodeExpr's guards: a non-Edge record, a
// non-PREDICATE-DECLARE codec, a missing Edge Id, a missing node count, a
// malformed label table, a malformed record stream, leftover record words and a
// node-count mismatch.
func TestDecodeExprErrors(t *testing.T) {
	rec := func(p ...core.Word) core.Record { return core.Record{Kind: core.KindEdge, Codec: 0x61, Payload: p} }
	if _, err := DecodeExpr(core.Record{Kind: core.KindValue}, testWire); err == nil {
		t.Fatal("non-Edge record: want error")
	}
	if _, err := DecodeExpr(core.Record{Kind: core.KindEdge, Codec: 0x60}, testWire); err == nil {
		t.Fatal("wrong codec: want error")
	}
	bad := []core.Record{rec(), rec(0x00E0, 1), rec(0x00E0, 1, 0, 1), rec(0x00E0, 1, 0, 0, 0)}
	for i, r := range bad {
		if _, err := DecodeExpr(r, testWire); err == nil {
			t.Fatalf("malformed payload %d: want error", i)
		}
	}
	ws, _ := EncodeExpr(ExprPacket{ID: 1, Cond: ast.Compare{Field: "f", Op: "is empty"}}, testWire)
	good, _, _ := core.DecodeRecord(ws)
	leftover := rec(append(append([]core.Word(nil), good.Payload...), 0)...)
	if _, err := DecodeExpr(leftover, testWire); err == nil {
		t.Fatal("leftover record words: want error")
	}
	mismatch := rec(append([]core.Word(nil), good.Payload...)...)
	mismatch.Payload[2]++
	if _, err := DecodeExpr(mismatch, testWire); err == nil {
		t.Fatal("node count mismatch: want error")
	}
}
