//ff:func feature=approve type=model control=iteration dimension=1
//ff:what TestThresholdSpecValidate — ThresholdSpec.Validate requires a positive Max
package approve

import "testing"

// TestThresholdSpecValidate covers both branches of ThresholdSpec.Validate:
// a non-positive Max (zero or negative) is an error, and a positive Max is
// valid.
func TestThresholdSpecValidate(t *testing.T) {
	cases := []struct {
		name    string
		max     float64
		wantErr bool
	}{
		{name: "negative max is invalid", max: -1, wantErr: true},
		{name: "zero max is invalid", max: 0, wantErr: true},
		{name: "positive max is valid", max: 10000, wantErr: false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			spec := &ThresholdSpec{Max: c.max}
			err := spec.Validate()
			if c.wantErr && err == nil {
				t.Error("Validate() = nil, want error")
			}
			if !c.wantErr && err != nil {
				t.Errorf("Validate() = %v, want nil", err)
			}
		})
	}
}
