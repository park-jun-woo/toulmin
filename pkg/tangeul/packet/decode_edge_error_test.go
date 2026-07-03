//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeEdgeErrors — DecodeEdge rejects a bad frame, truncated bodies, an unknown PropCode, a Block dst and a bad GATE16
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeEdgeErrors exercises every self-contained guard in DecodeEdge: a
// malformed frame header, a body shorter than the framed length, a body too
// short for PropCode+src, an unknown PropCode, a FuncRef dst with a truncated
// slot, a TID dst missing its word, a Block dst (provenance, not decoded in
// v0.1), a missing flag word, F_GATE with a missing gate slot, F_GATE with an
// invalid GATE16, and slot words that do not add up to the framed length.
func TestDecodeEdgeErrors(t *testing.T) {
	// A malformed frame header (bad prefix) fails at word.DecodeFrame.
	if _, _, err := DecodeEdge([]word.Word{0x0000}); err == nil {
		t.Fatalf("bad frame prefix: want error")
	}

	// A body shorter than the framed length.
	full, err := EncodeEdge(EdgePacket{Prop: codebook.PropAttacks, Src: 1, DstTID: 2, Span: word.Span{TID: 1}})
	if err != nil {
		t.Fatalf("EncodeEdge: %v", err)
	}
	if _, _, err := DecodeEdge(full[:len(full)-1]); err == nil {
		t.Fatalf("truncated body: want error")
	}

	// A body too short for PropCode+src (framed Len=1, one body word).
	short := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 1}), word.Word(codebook.PropAttacks))
	if _, _, err := DecodeEdge(short); err == nil {
		t.Fatalf("body shorter than PropCode+src: want error")
	}

	// An unknown PropCode fails at codebook.PropDstKind.
	unknown := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 2}), word.Word(0x0000), word.Word(1))
	if _, _, err := DecodeEdge(unknown); err == nil {
		t.Fatalf("unknown PropCode: want error")
	}

	// A FuncRef dst (do) whose body ends right after src — decodeFuncRef sees 0 words.
	funcShort := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 2}), word.Word(codebook.PropDo), word.Word(1))
	if _, _, err := DecodeEdge(funcShort); err == nil {
		t.Fatalf("truncated FuncRef dst: want error")
	}

	// A TID dst (attacks) whose body ends right after src — dst word missing.
	tidShort := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 2}), word.Word(codebook.PropAttacks), word.Word(1))
	if _, _, err := DecodeEdge(tidShort); err == nil {
		t.Fatalf("missing TID dst: want error")
	}

	// A Block dst (provenance) is not decoded in v0.1.
	block := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 2}), word.Word(codebook.PropProvenance), word.Word(1))
	if _, _, err := DecodeEdge(block); err == nil {
		t.Fatalf("Block dst: want error")
	}

	// A TID dst body that ends right after dst — the flag word is missing.
	noFlag := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 3}), word.Word(codebook.PropAttacks), word.Word(1), word.Word(2))
	if _, _, err := DecodeEdge(noFlag); err == nil {
		t.Fatalf("missing flag word: want error")
	}

	// F_GATE set but the body ends right after the flag word — the gate slot is missing.
	noGateSlot := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 4}), word.Word(codebook.PropAttacks), word.Word(1), word.Word(2), word.Word(codebook.EdgeFlagGate))
	if _, _, err := DecodeEdge(noGateSlot); err == nil {
		t.Fatalf("missing gate slot: want error")
	}

	// F_GATE set with an invalid GATE16 (reserved bits set).
	badGate := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 5}), word.Word(codebook.PropAttacks), word.Word(1), word.Word(2), word.Word(codebook.EdgeFlagGate), word.Word(0x0080))
	if _, _, err := DecodeEdge(badGate); err == nil {
		t.Fatalf("invalid GATE16: want error")
	}

	// Slot words do not add up to the framed length (one extra word before SPAN).
	span := word.EncodeSpan(word.Span{TID: 1})
	mismatch := append(word.EncodeFrame(word.Frame{Class: word.ClassEdge, Len: 8}), word.Word(codebook.PropAttacks), word.Word(1), word.Word(2), word.Word(0), word.Word(0xFFFF))
	mismatch = append(mismatch, span[:]...)
	if _, _, err := DecodeEdge(mismatch); err == nil {
		t.Fatalf("slot/length mismatch: want error")
	}
}
