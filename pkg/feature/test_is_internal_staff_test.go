//ff:func feature=feature type=rule control=iteration dimension=1
//ff:what TestIsInternalStaff — tests IsInternalStaff rule
package feature

import "testing"

func TestIsInternalStaff(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want bool
	}{
		{"internal", map[string]any{"internal": true}, true},
		{"not internal", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsInternalStaff(nil, &UserContext{Attributes: tt.attr}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
