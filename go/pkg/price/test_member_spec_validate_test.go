//ff:func feature=price type=engine control=sequence
//ff:what TestMemberSpec_Validate — tests validation of MemberSpec required Level field
package price

import "testing"

func TestMemberSpec_Validate(t *testing.T) {
	if err := (&MemberSpec{Level: ""}).Validate(); err == nil {
		t.Fatal("expected error for empty Level")
	}
	if err := (&MemberSpec{Level: "gold"}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
