//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestPackageName — tests packageName for valid chars, sanitized chars, empty input, and leading-digit branches
package gen

import "testing"

func TestPackageName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty input", "", "_"},
		{"valid alnum and underscore chars", "abc123_ABC", "abc123_ABC"},
		{"invalid chars replaced with underscore", "abc-def!ghi", "abc_def_ghi"},
		{"leading digit gets prefix", "123abc", "_123abc"},
		{"all invalid chars", "!!!", "___"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := packageName(tt.in)
			if got != tt.want {
				t.Errorf("packageName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
