//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLocaleCodeRoundtrip — LocaleCodeFromTag/LocaleTagFromCode invert for each locale and fall back to unknown
package codebook

import "testing"

// TestLocaleCodeRoundtrip checks that each authoring-locale tag maps to its code
// and back, and that an unknown tag or code degrades to LocaleUnknown / "".
func TestLocaleCodeRoundtrip(t *testing.T) {
	cases := []struct {
		tag  string
		code uint16
	}{
		{"en", LocaleEN},
		{"ko", LocaleKO},
		{"", LocaleUnknown},
	}
	for _, c := range cases {
		if got := LocaleCodeFromTag(c.tag); got != c.code {
			t.Errorf("LocaleCodeFromTag(%q) = %d, want %d", c.tag, got, c.code)
		}
		if got := LocaleTagFromCode(c.code); got != c.tag {
			t.Errorf("LocaleTagFromCode(%d) = %q, want %q", c.code, got, c.tag)
		}
	}
	if got := LocaleCodeFromTag("fr"); got != LocaleUnknown {
		t.Errorf("LocaleCodeFromTag(unknown) = %d, want LocaleUnknown", got)
	}
	if got := LocaleTagFromCode(0x7FFF); got != "" {
		t.Errorf("LocaleTagFromCode(unknown) = %q, want empty", got)
	}
}
