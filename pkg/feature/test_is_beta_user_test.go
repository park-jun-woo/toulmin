//ff:func feature=feature type=rule control=iteration dimension=1
//ff:what TestIsBetaUser — tests IsBetaUser rule
package feature

import "testing"

func TestIsBetaUser(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		want bool
	}{
		{"beta", map[string]any{"beta": true}, true},
		{"not beta", map[string]any{"beta": false}, false},
		{"no attr", map[string]any{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := IsBetaUser(nil, &UserContext{Attributes: tt.attr}, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
