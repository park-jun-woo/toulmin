//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunNode — tests runNode for RunOn nil/error/success and RunGraph nil/error/success branches
package toulmin

import (
	"errors"
	"strings"
	"testing"
)

func TestRunNode(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NoHandlerNoSubGraph", func(t *testing.T) {
			meta := &RuleMeta{Name: "n"}
			self := TraceEntry{Name: "n"}
			tr := Trace{}
			err := runNode(meta, self, tr, NewContext(), EvalOption{}, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}},
		{"HandlerError", func(t *testing.T) {
			meta := &RuleMeta{
				Name: "n",
				RunOn: func(self TraceEntry, tr Trace) error {
					return errors.New("boom")
				},
			}
			self := TraceEntry{Name: "n"}
			err := runNode(meta, self, Trace{}, NewContext(), EvalOption{}, 0)
			if err == nil {
				t.Fatal("expected handler error")
			}
			if !strings.Contains(err.Error(), "runOn") || !strings.Contains(err.Error(), "boom") {
				t.Errorf("expected wrapped runOn error, got %v", err)
			}
		}},
		{"HandlerSuccessAndSubGraphError", func(t *testing.T) {
			sub := NewGraph("sub")
			f1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			f2 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			a := sub.Rule(f1)
			b := sub.Counter(f2)
			a.Attacks(b)
			b.Attacks(a) // cycle -> runDepth returns error

			fired := false
			meta := &RuleMeta{
				Name: "n",
				RunOn: func(self TraceEntry, tr Trace) error {
					fired = true
					return nil
				},
				RunGraph: sub,
			}
			self := TraceEntry{Name: "n"}
			err := runNode(meta, self, Trace{}, NewContext(), EvalOption{}, 0)
			if !fired {
				t.Error("expected RunOn handler to fire")
			}
			if err == nil {
				t.Fatal("expected sub-graph error")
			}
			if !strings.Contains(err.Error(), "run ") || !strings.Contains(err.Error(), "→") {
				t.Errorf("expected wrapped sub-Run error, got %v", err)
			}
		}},
		{"SubGraphSuccess", func(t *testing.T) {
			sub := NewGraph("sub")
			subRan := false
			sub.Rule(WarrantA).RunOn(func(self TraceEntry, tr Trace) error {
				subRan = true
				return nil
			})

			meta := &RuleMeta{Name: "n", RunGraph: sub}
			self := TraceEntry{Name: "n"}
			err := runNode(meta, self, Trace{}, NewContext(), EvalOption{}, 0)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !subRan {
				t.Error("expected sub-graph to Run successfully")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
