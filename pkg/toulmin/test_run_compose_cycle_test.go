//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeCycle — a run cycle A→B→A is rejected by detectRunCycle before any execution
package toulmin

import (
	"strings"
	"testing"
)

func TestRunComposeCycle(t *testing.T) {
	ruleA := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	ruleB := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	a := NewGraph("A")
	b := NewGraph("B")
	a.Rule(ruleA).Run(b)
	b.Rule(ruleB).Run(a)

	results, trace, err := a.Run(NewContext())
	if err == nil {
		t.Fatal("expected run cycle error")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("error should mention cycle: %v", err)
	}
	if results != nil || trace != nil {
		t.Errorf("on cycle error Run must return nil results and trace, got results=%v trace=%v", results, trace)
	}
}
