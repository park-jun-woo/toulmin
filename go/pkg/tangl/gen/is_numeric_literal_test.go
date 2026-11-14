//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestIsNumericLiteral — tests isNumericLiteral for empty, sign-only, digit, decimal, double-dot, and non-numeric branches
package gen

import "testing"

func TestIsNumericLiteral(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"empty string", "", false},
		{"sign only", "+", false},
		{"negative sign only", "-", false},
		{"plain integer", "123", true},
		{"negative integer", "-123", true},
		{"positive integer", "+123", true},
		{"decimal", "1.5", true},
		{"leading dot only, no digit", ".", false},
		{"double dot invalid", "1.2.3", false},
		{"non-numeric char", "1a2", false},
		{"letters only", "abc", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNumericLiteral(tt.in)
			if got != tt.want {
				t.Errorf("isNumericLiteral(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
