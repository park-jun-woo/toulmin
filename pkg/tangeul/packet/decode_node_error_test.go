//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeNodeErrors — DecodeNode rejects a bad frame, a short header body and every missing flagged slot
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeNodeErrors exercises DecodeNode's self-contained guards: a
// malformed frame header, a body too short for TID+NodeType+header, a
// truncated F_FUNC FuncRef slot, a missing F_QUAL/F_CHECK/F_VERDICT slot, a
// malformed F_CLAIM block, a missing F_WITH count word and a truncated
// F_WITH TID list.
func TestDecodeNodeErrors(t *testing.T) {
	if _, _, err := DecodeNode([]word.Word{0x0000}); err == nil {
		t.Fatalf("bad frame prefix: want error")
	}

	short := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 2}), word.Word(1), word.Word(codebook.NodeField))
	if _, _, err := DecodeNode(short); err == nil {
		t.Fatalf("body shorter than TID+NodeType+header: want error")
	}

	funcTrunc := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 4}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagFunc), word.Word(2))
	if _, _, err := DecodeNode(funcTrunc); err == nil {
		t.Fatalf("truncated F_FUNC slot: want error")
	}

	qualMissing := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 3}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagQual))
	if _, _, err := DecodeNode(qualMissing); err == nil {
		t.Fatalf("missing F_QUAL slot: want error")
	}

	checkMissing := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 3}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagCheck))
	if _, _, err := DecodeNode(checkMissing); err == nil {
		t.Fatalf("missing F_CHECK slot: want error")
	}

	claimTrunc := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 4}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagClaim), word.Word(0x0011))
	if _, _, err := DecodeNode(claimTrunc); err == nil {
		t.Fatalf("malformed F_CLAIM block: want error")
	}

	verdictMissing := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 3}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagVerdict))
	if _, _, err := DecodeNode(verdictMissing); err == nil {
		t.Fatalf("missing F_VERDICT slot: want error")
	}

	withCountMissing := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 3}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagWith))
	if _, _, err := DecodeNode(withCountMissing); err == nil {
		t.Fatalf("missing F_WITH count: want error")
	}

	withListTrunc := append(word.EncodeFrame(word.Frame{Class: word.ClassNode, Len: 4}),
		word.Word(1), word.Word(codebook.NodeRule), nodeHeader(codebook.RoleGeneral, ObjTID, codebook.NodeFlagWith), word.Word(2))
	if _, _, err := DecodeNode(withListTrunc); err == nil {
		t.Fatalf("truncated F_WITH TID list: want error")
	}
}
