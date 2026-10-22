//ff:func feature=approve type=model control=iteration dimension=1
//ff:what TestApproverSpecValidate — ApproverSpec.Validate always reports success
package approve

import "testing"

// TestApproverSpecValidate covers the single branch of ApproverSpec.Validate:
// it unconditionally returns nil regardless of the receiver's Role and
// Level fields.
func TestApproverSpecValidate(t *testing.T) {
	cases := []struct {
		name string
		spec *ApproverSpec
	}{
		{name: "zero value", spec: &ApproverSpec{}},
		{name: "populated fields", spec: &ApproverSpec{Role: "manager", Level: 3}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.spec.Validate(); err != nil {
				t.Errorf("Validate() = %v, want nil", err)
			}
		})
	}
}
