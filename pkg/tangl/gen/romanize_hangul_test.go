//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what TestRomanizeHangul — tests romanizeHangul for full jamo, no-final, silent-initial, and non-syllable branches
package gen

import "testing"

func TestRomanizeHangul(t *testing.T) {
	tests := []struct {
		name string
		in   rune
		want string
	}{
		{"initial medial final", '출', "chul"},
		{"initial medial final geum", '금', "geum"},
		{"no final consonant", '가', "ga"},
		{"final ieung ng", '능', "neung"},
		{"silent initial ieung", '아', "a"},
		{"first syllable ga", '가', "ga"},
		{"last syllable hih", '힣', "hit"},
		{"non-syllable jamo dropped", 'ㄱ', ""},
		{"ascii not a syllable", 'a', ""},
		{"below block boundary", 0xABFF, ""},
		{"above block boundary", 0xD7A4, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := romanizeHangul(tt.in)
			if got != tt.want {
				t.Errorf("romanizeHangul(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
