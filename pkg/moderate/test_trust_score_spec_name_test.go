//ff:func feature=moderate type=engine control=sequence
//ff:what TestTrustScoreSpec_SpecName — tests SpecName returns the fixed spec name
package moderate

import "testing"

func TestTrustScoreSpec_SpecName(t *testing.T) {
	spec := &TrustScoreSpec{MinScore: 0.9}
	if got := spec.SpecName(); got != "TrustScoreSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "TrustScoreSpec")
	}
}
