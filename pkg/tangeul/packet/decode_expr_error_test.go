//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeExprErrors — DecodeExpr rejects a bad frame, truncated bodies, a malformed record stream, leftover words and a node-count mismatch
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestDecodeExprErrors exercises DecodeExpr's self-contained guards: a
// malformed frame header, a body shorter than the framed length, a body too
// short for a node count, a malformed label table, a malformed record stream,
// leftover unconsumed record words, and a node-count mismatch between the
// framed N and the decoded subtree.
func TestDecodeExprErrors(t *testing.T) {
	if _, _, err := DecodeExpr([]word.Word{0x0000}); err == nil {
		t.Fatalf("bad frame prefix: want error")
	}

	base, err := EncodeExpr(ExprPacket{Cond: ast.Compare{Field: "f", Op: "equals", Value: "0"}})
	if err != nil {
		t.Fatalf("EncodeExpr: %v", err)
	}
	if _, _, err := DecodeExpr(base[:len(base)-1]); err == nil {
		t.Fatalf("truncated body: want error")
	}

	empty := word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 0})
	if _, _, err := DecodeExpr(empty); err == nil {
		t.Fatalf("body shorter than a node count word: want error")
	}

	badLabels := append(word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 2}), word.Word(0), word.Word(5))
	if _, _, err := DecodeExpr(badLabels); err == nil {
		t.Fatalf("malformed label table: want error")
	}

	badRecord := append(word.EncodeFrame(word.Frame{Class: word.ClassPayload, Len: 2}), word.Word(0), word.Word(0))
	if _, _, err := DecodeExpr(badRecord); err == nil {
		t.Fatalf("malformed record stream: want error")
	}

	// Bump the header word's Len field by one and append an extra trailing word
	// so the record stream leaves one word unconsumed.
	leftover := append([]word.Word(nil), base...)
	leftover[0]++
	leftover = append(leftover, 0x0000)
	if _, _, err := DecodeExpr(leftover); err == nil {
		t.Fatalf("leftover record words: want error")
	}

	// Claim one extra node than the record stream actually decodes to.
	mismatch := append([]word.Word(nil), base...)
	mismatch[1]++
	if _, _, err := DecodeExpr(mismatch); err == nil {
		t.Fatalf("node count mismatch: want error")
	}
}
