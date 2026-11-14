//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseMarker — tests parseMarker for unordered/ordered markers, malformed input, and numeric overflow
package scan

import "testing"

func TestParseMarker(t *testing.T) {
	cases := []struct {
		name        string
		trimmed     string
		wantContent string
		wantOrdered bool
		wantNumber  int
		wantOK      bool
	}{
		{"empty", "", "", false, 0, false},
		{"unordered marker", "- x", "x", false, 0, true},
		{"dash without space", "-x", "", false, 0, false},
		{"non digit non dash", "abc", "", false, 0, false},
		{"ordered marker", "3. third", "third", true, 3, true},
		{"digits with no trailing dot-space", "12", "", false, 0, false},
		{"digits then non-dot", "12xy", "", false, 0, false},
		{"digits dot no space", "12.x", "", false, 0, false},
		{"digits dot space padded content", "42.  answer", "answer", true, 42, true},
		{"numeric overflow", "99999999999999999999999. x", "", false, 0, false},
	}
	for _, c := range cases {
		content, ordered, number, ok := parseMarker(c.trimmed)
		if content != c.wantContent || ordered != c.wantOrdered || number != c.wantNumber || ok != c.wantOK {
			t.Errorf("%s: parseMarker(%q) = (%q, %v, %d, %v), want (%q, %v, %d, %v)",
				c.name, c.trimmed, content, ordered, number, ok,
				c.wantContent, c.wantOrdered, c.wantNumber, c.wantOK)
		}
	}
}
