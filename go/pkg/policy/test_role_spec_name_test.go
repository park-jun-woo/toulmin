//ff:func feature=policy type=engine control=sequence
//ff:what TestRoleSpec_SpecName — verifies RoleSpec.SpecName returns fixed name
package policy

import "testing"

func TestRoleSpec_SpecName(t *testing.T) {
	if got := (&RoleSpec{}).SpecName(); got != "RoleSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "RoleSpec")
	}
}
