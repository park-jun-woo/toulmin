//ff:func feature=cli type=model control=iteration dimension=1
//ff:what TestDemoSpecSpecName — demoSpec.SpecName always returns the fixed identifier
package cli

import "testing"

// TestDemoSpecSpecName covers the single branch of demoSpec.SpecName: it
// unconditionally returns "demoSpec" regardless of the receiver's Value field.
func TestDemoSpecSpecName(t *testing.T) {
	cases := []struct {
		name string
		spec *demoSpec
	}{
		{name: "empty value", spec: &demoSpec{}},
		{name: "non-empty value", spec: &demoSpec{Value: "some spec"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.spec.SpecName(); got != "demoSpec" {
				t.Errorf("SpecName() = %q, want %q", got, "demoSpec")
			}
		})
	}
}
