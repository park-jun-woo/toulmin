//ff:func feature=approve type=rule control=iteration dimension=1
//ff:what TestIsUrgent — tests IsUrgent rule
package approve

import "testing"

func TestIsUrgent(t *testing.T) {
	tests := []struct {
		name string
		meta map[string]any
		want bool
	}{
		{"urgent", map[string]any{"urgent": true}, true},
		{"not urgent", map[string]any{}, false},
		{"nil meta", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ApprovalRequest{Metadata: tt.meta}
			got, _ := IsUrgent(req, nil, nil)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
