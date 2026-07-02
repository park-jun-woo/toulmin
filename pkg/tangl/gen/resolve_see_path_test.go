//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestResolveSeePath — tests resolveSeePath for tangl-prefixed and pass-through branches
package gen

import "testing"

func TestResolveSeePath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"tangl prefix rewritten to full module path", "tangl/orders", "github.com/park-jun-woo/toulmin/pkg/tangl/orders"},
		{"non-tangl path passed through unchanged", "example.com/pkg/orders", "example.com/pkg/orders"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveSeePath(tt.in)
			if got != tt.want {
				t.Errorf("resolveSeePath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
