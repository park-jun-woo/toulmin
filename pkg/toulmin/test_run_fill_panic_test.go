//ff:func feature=engine type=engine control=sequence
//ff:what TestRunFillPanic — a node that panics only during the full pass aborts Run with an error
package toulmin

import (
	"strings"
	"testing"
)

func TestRunFillPanic(t *testing.T) {
	panicCounter := func(ctx Context, specs Specs) (bool, any) {
		panic("fill panic")
	}
	g := NewGraph("fillpanic")
	// InactiveR returns false, so lazy evaluation never recurses into its attacker.
	w := g.Rule(InactiveR)
	c := g.Counter(panicCounter)
	c.Attacks(w)

	results, trace, err := g.Run(NewContext())
	if err == nil {
		t.Fatal("expected error from node panicking during the full pass")
	}
	if !strings.Contains(err.Error(), "panicked") {
		t.Errorf("expected panic error, got: %v", err)
	}
	if results != nil || trace != nil {
		t.Errorf("on full-pass error Run must return nil results and nil trace, got results=%v trace=%v", results, trace)
	}
}
