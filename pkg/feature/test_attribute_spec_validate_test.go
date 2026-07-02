//ff:func feature=feature type=engine control=iteration dimension=1
//ff:what TestAttributeSpec_Validate — tests validation of AttributeSpec key presence
package feature

import "testing"

func TestAttributeSpec_Validate(t *testing.T) {
	cases := []struct {
		name    string
		spec    *AttributeSpec
		wantErr bool
	}{
		{
			name:    "empty key",
			spec:    &AttributeSpec{Key: "", Value: true},
			wantErr: true,
		},
		{
			name:    "non-empty key",
			spec:    &AttributeSpec{Key: "beta", Value: true},
			wantErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.spec.Validate()
			if c.wantErr && err == nil {
				t.Fatalf("Validate() = nil, want error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("Validate() = %v, want nil", err)
			}
		})
	}
}
