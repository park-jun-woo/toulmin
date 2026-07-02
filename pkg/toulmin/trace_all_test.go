//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestTraceAll — tests Trace.All returns the underlying node slice
package toulmin

import "testing"

func TestTraceAll(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NonEmpty", func(t *testing.T) {
			nodes := []TraceEntry{{Name: "a"}, {Name: "b"}}
			tr := Trace{nodes: nodes}
			got := tr.All()
			if len(got) != 2 || got[0].Name != "a" || got[1].Name != "b" {
				t.Fatalf("expected %v, got %v", nodes, got)
			}
		}},
		{"Empty", func(t *testing.T) {
			tr := Trace{}
			if got := tr.All(); got != nil {
				t.Fatalf("expected nil, got %v", got)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
