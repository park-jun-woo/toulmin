//ff:func feature=policy type=engine control=sequence
//ff:what TestRoleSpec_Validate — tests validation of RoleSpec required Role field
package policy

import "testing"

func TestRoleSpec_Validate(t *testing.T) {
	if err := (&RoleSpec{Role: ""}).Validate(); err == nil {
		t.Fatal("expected error for empty Role")
	}
	if err := (&RoleSpec{Role: "admin"}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
