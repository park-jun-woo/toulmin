//ff:func feature=cli type=model control=iteration dimension=1
//ff:what TestDemoSpecValidate — demoSpec.Validate always reports success
package cli

import "testing"

// TestDemoSpecValidate covers the single branch of demoSpec.Validate: it
// unconditionally returns nil regardless of the receiver's Value field.
func TestDemoSpecValidate(t *testing.T) {
	cases := []struct {
		name string
		spec *demoSpec
	}{
		{name: "empty value", spec: &demoSpec{}},
		{name: "non-empty value", spec: &demoSpec{Value: "some spec"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.spec.Validate(); err != nil {
				t.Errorf("Validate() = %v, want nil", err)
			}
		})
	}
}
