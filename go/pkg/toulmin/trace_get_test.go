//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestTraceGet — tests Trace.Get for found and not-found branches
package toulmin

import "testing"

func TestTraceGet(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Found", func(t *testing.T) {
			tr := Trace{nodes: []TraceEntry{{Name: "a"}, {Name: "b"}}}
			entry, ok := tr.Get("b")
			if !ok {
				t.Fatalf("expected ok=true")
			}
			if entry.Name != "b" {
				t.Fatalf("expected entry Name=b, got %q", entry.Name)
			}
		}},
		{"NotFound", func(t *testing.T) {
			tr := Trace{nodes: []TraceEntry{{Name: "a"}}}
			entry, ok := tr.Get("missing")
			if ok {
				t.Fatalf("expected ok=false")
			}
			if entry.Name != "" || entry.Specs != nil || entry.Evidence != nil {
				t.Fatalf("expected zero-value TraceEntry, got %+v", entry)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
