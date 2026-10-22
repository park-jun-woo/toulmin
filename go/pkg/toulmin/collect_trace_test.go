//ff:func feature=engine type=util control=iteration dimension=1
//ff:what TestCollectTrace — tests collectTrace for empty and populated trace branches
package toulmin

import "testing"

func TestCollectTrace(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Empty", func(t *testing.T) {
			out := collectTrace(nil)
			if len(out) != 0 {
				t.Fatalf("expected empty slice, got %v", out)
			}
		}},
		{"Populated", func(t *testing.T) {
			in := []TraceEntry{
				{Name: "pkg.Foo", Verdict: 0.5},
				{Name: "Bar", Verdict: -0.5},
			}
			out := collectTrace(in)
			if len(out) != 2 {
				t.Fatalf("expected 2 entries, got %d", len(out))
			}
			if out[0].Name != "Foo" {
				t.Errorf("out[0].Name = %q, want %q", out[0].Name, "Foo")
			}
			if out[1].Name != "Bar" {
				t.Errorf("out[1].Name = %q, want %q", out[1].Name, "Bar")
			}
			// Ensure the input slice was not mutated in place.
			if in[0].Name != "pkg.Foo" {
				t.Errorf("input mutated: in[0].Name = %q", in[0].Name)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
