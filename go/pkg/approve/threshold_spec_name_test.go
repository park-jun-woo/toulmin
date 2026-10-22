//ff:func feature=approve type=model control=iteration dimension=1
//ff:what TestThresholdSpecSpecName — ThresholdSpec.SpecName always returns the fixed identifier
package approve

import "testing"

// TestThresholdSpecSpecName covers the single branch of
// ThresholdSpec.SpecName: it unconditionally returns "ThresholdSpec"
// regardless of the receiver's Max field.
func TestThresholdSpecSpecName(t *testing.T) {
	cases := []struct {
		name string
		spec *ThresholdSpec
	}{
		{name: "zero value", spec: &ThresholdSpec{}},
		{name: "populated field", spec: &ThresholdSpec{Max: 10000}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.spec.SpecName(); got != "ThresholdSpec" {
				t.Errorf("SpecName() = %q, want %q", got, "ThresholdSpec")
			}
		})
	}
}
