//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeNodeErrors — EncodeNode rejects nil Qual/Check/Claim/Verdict fields and an invalid Claim
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEncodeNodeErrors exercises EncodeNode's remaining nil-field guards
// (F_QUAL, F_CHECK, F_CLAIM, F_VERDICT) and the propagated encodeClaim error
// when F_CLAIM is set with a Claim of unknown ObjKind.
func TestEncodeNodeErrors(t *testing.T) {
	base := NodePacket{TID: 1, Type: codebook.NodeRule, Role: codebook.RoleGeneral, Span: word.Span{TID: 1}}

	qual := base
	qual.Flags = codebook.NodeFlagQual
	if _, err := EncodeNode(qual); err == nil {
		t.Fatalf("F_QUAL without Qual: want error")
	}

	check := base
	check.Flags = codebook.NodeFlagCheck
	if _, err := EncodeNode(check); err == nil {
		t.Fatalf("F_CHECK without Check: want error")
	}

	claim := base
	claim.Flags = codebook.NodeFlagClaim
	if _, err := EncodeNode(claim); err == nil {
		t.Fatalf("F_CLAIM without Claim: want error")
	}

	verdict := base
	verdict.Flags = codebook.NodeFlagVerdict
	if _, err := EncodeNode(verdict); err == nil {
		t.Fatalf("F_VERDICT without Verdict: want error")
	}

	badClaim := base
	badClaim.Flags = codebook.NodeFlagClaim
	badClaim.Claim = &Claim{Kind: ObjKind(3)}
	if _, err := EncodeNode(badClaim); err == nil {
		t.Fatalf("F_CLAIM with unknown ObjKind: want error")
	}
}
