//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestGoIdent — tests goIdent for multi-word, empty-after-sanitize, all-empty, and leading-digit branches
package gen

import "testing"

func TestGoIdent(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"single word", "order", "order"},
		{"multi word camel", "order received", "orderReceived"},
		{"punctuation-only word skipped", "order !!! received", "orderReceived"},
		{"all words empty after sanitize", "!!! ###", "_"},
		{"empty input", "", "_"},
		{"leading digit gets underscore prefix", "123abc", "_123abc"},
		{"hangul multi word romanized camel", "출금 가능", "chulgeumGaneung"},
		{"hangul deposit distinct from withdrawal", "입금 가능", "ipgeumGaneung"},
		{"hangul balance sufficient", "잔액 충분", "janaekChungbun"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goIdent(tt.in)
			if got != tt.want {
				t.Errorf("goIdent(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
