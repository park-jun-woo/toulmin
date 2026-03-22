//ff:func feature=feature type=rule control=iteration dimension=1
//ff:what TestHasAttribute — tests HasAttribute rule
package feature

import "testing"

func TestHasAttribute(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		pair [2]any
		want bool
	}{
		{"match", map[string]any{"plan": "pro"}, [2]any{"plan", "pro"}, true},
		{"mismatch", map[string]any{"plan": "free"}, [2]any{"plan", "pro"}, false},
		{"missing", map[string]any{}, [2]any{"plan", "pro"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := HasAttribute(nil, &UserContext{Attributes: tt.attr}, tt.pair)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
