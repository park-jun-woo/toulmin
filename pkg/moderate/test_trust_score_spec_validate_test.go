//ff:func feature=moderate type=engine control=sequence
//ff:what TestTrustScoreSpec_Validate — tests validation of TrustScoreSpec MinScore bound [0,1]
package moderate

import "testing"

func TestTrustScoreSpec_Validate(t *testing.T) {
	if err := (&TrustScoreSpec{MinScore: -0.1}).Validate(); err == nil {
		t.Fatal("expected error for MinScore below 0")
	}
	if err := (&TrustScoreSpec{MinScore: 1.1}).Validate(); err == nil {
		t.Fatal("expected error for MinScore above 1")
	}
	if err := (&TrustScoreSpec{MinScore: 0.5}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
