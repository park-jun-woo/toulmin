//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseUnitDecimalKo — KO weight literal: value with attached 의 remainder, missing digits, bad float, >4 decimals, out of range
package ko

import (
	"strings"
	"testing"
)

func TestParseUnitDecimalKo(t *testing.T) {
	v, rest, err := parseUnitDecimalKo(" 0.50의", "세기", 1, "t.md")
	if err != nil || v != 0.5 || rest != "의" {
		t.Fatalf("got (%v, %q, %v), want (0.5, \"의\", nil)", v, rest, err)
	}
	cases := []struct {
		in, want string
	}{
		{"의", "expected float after '세기'"},
		{"--1", "invalid 세기 value"},
		{"0.12345의", "more than 4 decimal places"},
		{"2의", "out of range"},
	}
	for _, tc := range cases {
		_, _, err := parseUnitDecimalKo(tc.in, "세기", 4, "t.md")
		if err == nil || !strings.Contains(err.Error(), tc.want) || !strings.HasPrefix(err.Error(), "t.md:4:") {
			t.Errorf("parseUnitDecimalKo(%q) err = %v, want t.md:4: ... %q", tc.in, err, tc.want)
		}
	}
}
