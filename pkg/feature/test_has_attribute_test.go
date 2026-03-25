//ff:func feature=feature type=rule control=iteration dimension=1
//ff:what TestHasAttribute — tests HasAttribute rule
package feature

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestHasAttribute(t *testing.T) {
	tests := []struct {
		name string
		attr map[string]any
		key  string
		val  any
		want bool
	}{
		{"match", map[string]any{"plan": "pro"}, "plan", "pro", true},
		{"mismatch", map[string]any{"plan": "free"}, "plan", "pro", false},
		{"missing", map[string]any{}, "plan", "pro", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := toulmin.NewContext()
			ctx.Set("attributes", tt.attr)
			got, _ := HasAttribute(ctx, toulmin.Specs{&AttributeSpec{Key: tt.key, Value: tt.val}})
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
