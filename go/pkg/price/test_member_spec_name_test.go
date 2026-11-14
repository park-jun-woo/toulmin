//ff:func feature=price type=engine control=sequence
//ff:what TestMemberSpec_SpecName — verifies MemberSpec.SpecName returns fixed name
package price

import "testing"

func TestMemberSpec_SpecName(t *testing.T) {
	if got := (&MemberSpec{}).SpecName(); got != "MemberSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "MemberSpec")
	}
}
