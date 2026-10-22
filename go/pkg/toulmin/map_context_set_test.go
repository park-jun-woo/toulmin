//ff:func feature=engine type=model control=sequence
//ff:what TestMapContextSet — tests Set stores value retrievable via Get
package toulmin

import "testing"

func TestMapContextSet(t *testing.T) {
	c := NewContext()
	c.Set("k", "v")
	v, ok := c.Get("k")
	if !ok || v != "v" {
		t.Fatalf("expected v=%q ok=true, got v=%v ok=%v", "v", v, ok)
	}
	c.Set("k", "v2")
	v, ok = c.Get("k")
	if !ok || v != "v2" {
		t.Fatalf("expected overwritten value v2, got v=%v ok=%v", v, ok)
	}
}
