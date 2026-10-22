//ff:func feature=approve type=model control=iteration dimension=1
//ff:what TestApproverSpecSpecName — ApproverSpec.SpecName always returns the fixed identifier
package approve

import "testing"

// TestApproverSpecSpecName covers the single branch of ApproverSpec.SpecName:
// it unconditionally returns "ApproverSpec" regardless of the receiver's
// Role and Level fields.
func TestApproverSpecSpecName(t *testing.T) {
	cases := []struct {
		name string
		spec *ApproverSpec
	}{
		{name: "zero value", spec: &ApproverSpec{}},
		{name: "populated fields", spec: &ApproverSpec{Role: "manager", Level: 3}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.spec.SpecName(); got != "ApproverSpec" {
				t.Errorf("SpecName() = %q, want %q", got, "ApproverSpec")
			}
		})
	}
}
