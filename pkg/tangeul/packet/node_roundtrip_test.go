//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestNodeRoundTrip — Node packets across flag combinations round-trip and nil flagged fields error
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestNodeRoundTrip exercises Node packets across flag combinations — bare, a
// full RuleNode (F_FUNC|F_CLAIM|F_VERDICT), a qualifier+check node, a `with`
// clause node, and header bits 10-9 Level fail/review nodes — asserting
// decode(encode(p)) == p and that a flag whose field is nil is rejected.
func TestNodeRoundTrip(t *testing.T) {
	qual := word.VFix(8000)
	verdict := word.VFix(0)
	check := TID(0x0020)
	packets := []NodePacket{
		{TID: 0x0011, Type: codebook.NodeField, Role: codebook.RoleGeneral, Span: word.Span{TID: 1, Off: 2, Len: 3}},
		{
			TID: 0x0013, Type: codebook.NodeRule, Role: codebook.RoleGeneral,
			Flags:   codebook.NodeFlagFunc | codebook.NodeFlagClaim | codebook.NodeFlagVerdict,
			Func:    &FuncRef{Alias: 0x0002, Verb: 0x7400},
			Claim:   &Claim{Subject: 0x0011, Predicate: 0x7400, Kind: ObjTID, ObjTID: 0x0012},
			Verdict: &verdict,
			Span:    word.Span{TID: 1, Off: 40, Len: 55},
		},
		{
			TID: 0x0014, Type: codebook.NodeRule, Role: codebook.RoleCounter,
			Flags: codebook.NodeFlagQual | codebook.NodeFlagCheck,
			Qual:  &qual, Check: &check,
			Span: word.Span{TID: 1, Off: 60, Len: 40},
		},
		{
			TID: 0x0015, Type: codebook.NodeRule, Role: codebook.RoleGeneral,
			Flags: codebook.NodeFlagWith,
			With:  []TID{0x0007, 0x0008},
			Span:  word.Span{TID: 1, Off: 5, Len: 5},
		},
		{
			TID: 0x0016, Type: codebook.NodeRule, Role: codebook.RoleCounter,
			Level: codebook.LevelFail,
			Span:  word.Span{TID: 1, Off: 7, Len: 4},
		},
		{
			TID: 0x0017, Type: codebook.NodeRule, Role: codebook.RoleExcept,
			Level: codebook.LevelReview,
			Flags: codebook.NodeFlagQual, Qual: &qual,
			Span: word.Span{TID: 1, Off: 11, Len: 4},
		},
	}
	for i, in := range packets {
		ws, err := EncodeNode(in)
		if err != nil {
			t.Fatalf("packet %d: EncodeNode: %v", i, err)
		}
		got, n, err := DecodeNode(ws)
		if err != nil || n != len(ws) || !reflect.DeepEqual(got, in) {
			t.Fatalf("packet %d: DecodeNode = %+v, %d, %v; want %+v", i, got, n, err, in)
		}
	}
	if _, err := EncodeNode(NodePacket{Type: codebook.NodeRule, Flags: codebook.NodeFlagFunc}); err == nil {
		t.Fatalf("F_FUNC without Func: want error")
	}
}
