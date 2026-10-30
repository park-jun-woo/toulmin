//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseUnitDecimal — weight literal: value and remainder, missing digits, bad float, >4 decimals, out of range
package parser

import (
	"strings"
	"testing"
)

func TestParseUnitDecimal(t *testing.T) {
	v, rest, err := parseUnitDecimal("  0.50, level fail", "force", "t.md", 1)
	if err != nil || v != 0.5 || rest != ", level fail" {
		t.Fatalf("got (%v, %q, %v), want (0.5, \", level fail\", nil)", v, rest, err)
	}
	v, rest, err = parseUnitDecimal("1.0", "force", "t.md", 1)
	if err != nil || v != 1 || rest != "" {
		t.Fatalf("got (%v, %q, %v), want (1, \"\", nil)", v, rest, err)
	}
	cases := []struct {
		in, want string
	}{
		{"abc", "expected float after 'force'"},
		{"--1", "invalid force value \"--1\""},
		{"0.12345", "more than 4 decimal places"},
		{"1.5", "out of range"},
		{"-0.5", "out of range"},
	}
	for _, tc := range cases {
		_, _, err := parseUnitDecimal(tc.in, "force", "t.md", 7)
		if err == nil || !strings.Contains(err.Error(), tc.want) || !strings.HasPrefix(err.Error(), "t.md:7:") {
			t.Errorf("parseUnitDecimal(%q) err = %v, want t.md:7: ... %q", tc.in, err, tc.want)
		}
	}
}
