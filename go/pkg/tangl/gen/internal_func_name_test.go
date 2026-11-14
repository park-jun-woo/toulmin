//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestInternalFuncName — tests internalFuncName for underscore-fallback and identifier-suffix branches
package gen

import "testing"

func TestInternalFuncName(t *testing.T) {
	tests := []struct {
		name string
		kind string
		seed string
		idx  int
		want string
	}{
		{"empty seed falls back to kind+idx", "check", "", 3, "check3"},
		{"punctuation-only seed falls back to kind+idx", "check", "!!!", 1, "check1"},
		{"valid seed appends exported ident", "check", "order received", 2, "checkOrderReceived2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := internalFuncName(tt.kind, tt.seed, tt.idx)
			if got != tt.want {
				t.Errorf("internalFuncName(%q, %q, %d) = %q, want %q", tt.kind, tt.seed, tt.idx, got, tt.want)
			}
		})
	}
}
