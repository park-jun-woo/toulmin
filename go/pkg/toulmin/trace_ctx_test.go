//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestTraceCtx — tests Trace.Ctx returns the underlying context
package toulmin

import "testing"

func TestTraceCtx(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NonNil", func(t *testing.T) {
			ctx := NewContext()
			ctx.Set("k", "v")
			tr := Trace{ctx: ctx}
			got := tr.Ctx()
			if got != Context(ctx) {
				t.Fatalf("expected same context instance")
			}
			if v, ok := got.Get("k"); !ok || v != "v" {
				t.Fatalf("expected k=v, got v=%v ok=%v", v, ok)
			}
		}},
		{"Nil", func(t *testing.T) {
			tr := Trace{}
			if got := tr.Ctx(); got != nil {
				t.Fatalf("expected nil context, got %v", got)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
