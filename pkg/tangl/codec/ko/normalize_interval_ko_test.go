//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestNormalizeIntervalKo — 초/분/시간 unit conversion, clock passthrough, empty rejection
package ko

import "testing"

func TestNormalizeIntervalKo(t *testing.T) {
	cases := []struct {
		raw, want string
	}{
		{"30초", "30s"},
		{"1초", "1s"},
		{"5분", "5m"},
		{"1시간", "1h"},
		{"7:00", "7:00"},
	}
	for _, c := range cases {
		got, ok := normalizeIntervalKo(c.raw)
		if !ok || got != c.want {
			t.Errorf("normalizeIntervalKo(%q) = %q ok=%v, want %q", c.raw, got, ok, c.want)
		}
	}
	if _, ok := normalizeIntervalKo(""); ok {
		t.Error("expected ok=false for an empty interval")
	}
	if _, ok := normalizeIntervalKo("초"); ok {
		t.Error("expected ok=false for a unit with no number")
	}
}
