//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestSanitizeWord — tests sanitizeWord for alnum-kept and non-alnum-dropped branches
package gen

import "testing"

func TestSanitizeWord(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty input", "", ""},
		{"all alnum kept", "abcXYZ123", "abcXYZ123"},
		{"punctuation dropped", "a-b_c!d", "abcd"},
		{"all punctuation dropped to empty", "!!!", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeWord(tt.in)
			if got != tt.want {
				t.Errorf("sanitizeWord(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
