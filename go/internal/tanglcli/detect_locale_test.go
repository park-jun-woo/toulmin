//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestDetectLocale — the first `## tangl:` header keyword pins en/ko, and a missing/unknown keyword errors
package tanglcli

import (
	"strings"
	"testing"
)

// TestDetectLocale covers detectLocale's branches: an EN section keyword yields
// "en", a KO keyword yields "ko" (even when preceded by prose), an unknown
// keyword in the first header errors, a source with no `## tangl:` header
// errors, and a scanner failure (a single line exceeding the max token size)
// surfaces sc.Err() directly. Each case is checked by checkDetectLocale.
func TestDetectLocale(t *testing.T) {
	cases := []struct {
		name       string
		src        string
		wantLocale string
		wantErr    string
	}{
		{"english subject", "# Title\n\n## tangl:Subject\n- x\n", "en", ""},
		{"korean cases after prose", "설명 문단\n\n## tangl:사안\n- x\n", "ko", ""},
		{"unknown keyword", "## tangl:Bogus\n- x\n", "", "unknown section keyword"},
		{"no header", "just prose, nothing here\n", "", "no `## tangl:` section header"},
		{"scanner error on oversized line", strings.Repeat("a", 2<<20), "", "token too long"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { checkDetectLocale(t, c.src, c.wantLocale, c.wantErr) })
	}
}
