//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what TestCertaintyOpKo — each threshold token maps to its canonical Op, unknown rejected
package ko

import "testing"

func TestCertaintyOpKo(t *testing.T) {
	want := map[string]string{
		"이상이면": "at least",
		"초과면":  "above",
		"미만이면": "less than",
		"이하면":  "at most",
	}
	for token, op := range want {
		got, ok := certaintyOpKo(token)
		if !ok || got != op {
			t.Errorf("certaintyOpKo(%q) = %q ok=%v, want %q", token, got, ok, op)
		}
	}
	if _, ok := certaintyOpKo("근처면"); ok {
		t.Error("expected ok=false for an unknown token")
	}
}
