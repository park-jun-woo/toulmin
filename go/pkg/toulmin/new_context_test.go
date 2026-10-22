//ff:func feature=engine type=engine control=sequence
//ff:what TestNewContext — tests NewContext returns a usable empty MapContext
package toulmin

import "testing"

func TestNewContext(t *testing.T) {
	c := NewContext()
	if c == nil {
		t.Fatalf("expected non-nil context")
	}
	if _, ok := c.Get("anything"); ok {
		t.Fatalf("expected empty context")
	}
}
