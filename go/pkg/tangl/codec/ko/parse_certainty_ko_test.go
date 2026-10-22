//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestParseCertaintyKo — the four threshold ops map to canonical Op, plus error paths
package ko

import "testing"

func TestParseCertaintyKo(t *testing.T) {
	cases := []struct {
		text string
		op   string
	}{
		{"확신 75% 이상이면", "at least"},
		{"확신 90% 초과면", "above"},
		{"확신 20% 미만이면", "less than"},
		{"확신 50% 이하면", "at most"},
	}
	for _, c := range cases {
		cert, err := parseCertaintyKo(c.text, 1, "t.md")
		if err != nil {
			t.Fatalf("%q: %v", c.text, err)
		}
		if cert.Op != c.op {
			t.Errorf("%q: Op=%q, want %q", c.text, cert.Op, c.op)
		}
	}

	if _, err := parseCertaintyKo("확신 없음", 1, "t.md"); err == nil {
		t.Error("expected error for a missing percent")
	}
	if _, err := parseCertaintyKo("확신 75 이상이면", 1, "t.md"); err == nil {
		t.Error("expected error for a missing % sign")
	}
	if _, err := parseCertaintyKo("확신 75% 대충이면", 1, "t.md"); err == nil {
		t.Error("expected error for an unknown threshold operator")
	}
	if _, err := parseCertaintyKo("믿음 75% 이상이면", 1, "t.md"); err == nil {
		t.Error("expected error when the 확신 keyword is absent")
	}
	if _, err := parseCertaintyKo("확신 99999999999999999999% 이상이면", 1, "t.md"); err == nil {
		t.Error("expected error when the percent digit run overflows int (strconv.Atoi failure)")
	}
}
