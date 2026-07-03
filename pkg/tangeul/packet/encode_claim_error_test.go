//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeClaimErrors — encodeClaim rejects an unknown ObjKind
package packet

import "testing"

// TestEncodeClaimErrors checks encodeClaim rejects a Claim with an unknown
// ObjKind rather than silently defaulting an object slot.
func TestEncodeClaimErrors(t *testing.T) {
	if _, err := encodeClaim(Claim{Kind: ObjKind(3)}); err == nil {
		t.Fatalf("unknown ObjKind: want error")
	}
}
