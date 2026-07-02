//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestSplitLines — tests splitLines for CRLF normalization and splitting
package parser

import "testing"

func TestSplitLines(t *testing.T) {
	lines := splitLines("a\r\nb\nc")
	want := []string{"a", "b", "c"}
	if len(lines) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(lines), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line %d = %q, want %q", i, lines[i], w)
		}
	}
}
