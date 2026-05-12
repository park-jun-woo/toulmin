//ff:func feature=engine type=engine control=sequence
//ff:what TestMultipleWith — tests chained With calls accumulate specs
package toulmin

import (
	"testing"
)

func TestMultipleWith(t *testing.T) {
	fn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Rule(fn).With(&testSpec{Value: "a"}).With(&testSpec{Value: "b"})
	results, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if len(results[0].Trace[0].Specs) != 2 {
		t.Errorf("expected 2 specs, got %d", len(results[0].Trace[0].Specs))
	}
}
