//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestRomanizeRune — tests romanizeRune for Hangul-romanized and non-Hangul-passthrough branches
package gen

import "testing"

func TestRomanizeRune(t *testing.T) {
	tests := []struct {
		name string
		in   rune
		want string
	}{
		{"hangul syllable romanized", '출', "chul"},
		{"ascii letter passes through", 'a', "a"},
		{"digit passes through", '7', "7"},
		{"space passes through", ' ', " "},
		{"punctuation passes through", '-', "-"},
		{"standalone jamo passes through", 'ㄱ', "ㄱ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := romanizeRune(tt.in)
			if got != tt.want {
				t.Errorf("romanizeRune(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
