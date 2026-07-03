//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodePrefixItem — EXPR vs Text decoding, each on both a valid and a truncated packet
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestDecodePrefixItem covers both isExpr branches of decodePrefixItem and,
// within each, the decode-success and decode-error paths of the underlying
// packet.DecodeExpr / packet.DecodeText call.
func TestDecodePrefixItem(t *testing.T) {
	t.Run("exprOK", func(t *testing.T) {
		d := newDec()
		ws, err := packet.EncodeExpr(packet.ExprPacket{Cond: ast.Compare{Field: "level", Op: "is at least", Value: "5"}})
		if err != nil {
			t.Fatalf("EncodeExpr: %v", err)
		}
		got, err := decodePrefixItem(d, rawPacket{Class: word.ClassPayload, Words: ws}, true)
		if err != nil {
			t.Fatalf("decodePrefixItem: %v", err)
		}
		if got == nil {
			t.Errorf("got nil Expr, want the decoded Compare")
		}
	})

	t.Run("exprError", func(t *testing.T) {
		d := newDec()
		ws, err := packet.EncodeExpr(packet.ExprPacket{Cond: ast.Compare{Field: "level", Op: "is at least", Value: "5"}})
		if err != nil {
			t.Fatalf("EncodeExpr: %v", err)
		}
		bad := ws[:len(ws)-1]
		got, err := decodePrefixItem(d, rawPacket{Class: word.ClassPayload, Words: bad}, true)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if got != nil {
			t.Errorf("got %v, want nil on error", got)
		}
	})

	t.Run("textOK", func(t *testing.T) {
		d := newDec()
		before := d.textNext
		ws := packet.EncodeText(packet.TextPacket{Bytes: []byte("hello")})
		got, err := decodePrefixItem(d, rawPacket{Class: word.ClassPayload, Words: ws}, false)
		if err != nil {
			t.Fatalf("decodePrefixItem: %v", err)
		}
		if got != nil {
			t.Errorf("got %v, want nil Expr for a Text item", got)
		}
		if d.texts[before] != "hello" {
			t.Errorf("d.texts[%d] = %q, want %q", before, d.texts[before], "hello")
		}
		if d.textNext != before+1 {
			t.Errorf("textNext = %d, want %d", d.textNext, before+1)
		}
	})

	t.Run("textError", func(t *testing.T) {
		d := newDec()
		ws := packet.EncodeText(packet.TextPacket{Bytes: []byte("hello")})
		bad := ws[:len(ws)-1]
		got, err := decodePrefixItem(d, rawPacket{Class: word.ClassPayload, Words: bad}, false)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if got != nil {
			t.Errorf("got %v, want nil on error", got)
		}
	})
}
