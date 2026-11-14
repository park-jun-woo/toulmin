//ff:func feature=engine type=model control=iteration dimension=1
//ff:what TestMapContextGet — tests Get for found and not-found branches
package toulmin

import "testing"

func TestMapContextGet(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Found", func(t *testing.T) {
			c := NewContext()
			c.Set("k", 42)
			v, ok := c.Get("k")
			if !ok {
				t.Fatalf("expected ok=true")
			}
			if v != 42 {
				t.Fatalf("expected 42, got %v", v)
			}
		}},
		{"NotFound", func(t *testing.T) {
			c := NewContext()
			v, ok := c.Get("missing")
			if ok {
				t.Fatalf("expected ok=false")
			}
			if v != nil {
				t.Fatalf("expected nil, got %v", v)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
