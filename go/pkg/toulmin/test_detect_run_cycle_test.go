//ff:func feature=engine type=validator control=iteration dimension=1
//ff:what TestDetectRunCycle — detectRunCycle accepts DAGs (incl. diamonds) and rejects cycles
package toulmin

import (
	"strings"
	"testing"
)

// TestDetectRunCycle drives detectRunCycle directly over hand-built graph-of-graphs:
// an acyclic chain (nil), a self-loop and an A→B→A loop (error), and a diamond DAG
// where one sub-graph is reached by two paths (nil — done color short-circuits).
func TestDetectRunCycle(t *testing.T) {
	// Distinct function values: rules are keyed by function identity, so two rules in the
	// same graph must be different closures (reuse across graphs is fine).
	f1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	f2 := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	// acyclic: root → child (and a rule with no RunGraph to hit the `sub == nil` continue).
	child := NewGraph("child")
	child.Rule(f1)
	acyclic := NewGraph("acyclic")
	acyclic.Rule(f1) // no RunGraph -> continue branch
	acyclic.Rule(f2).Run(child)

	// self-loop: root → root.
	selfLoop := NewGraph("self")
	selfLoop.Rule(f1)
	selfLoop.rules[0].RunGraph = selfLoop

	// A → B → A loop.
	ga := NewGraph("A")
	gb := NewGraph("B")
	ga.Rule(f1).Run(gb)
	gb.Rule(f1).Run(ga)

	// diamond: parent → {shared, shared} (shared reached twice, no cycle).
	shared := NewGraph("shared")
	shared.Rule(f1)
	diamond := NewGraph("diamond")
	diamond.Rule(f1).Run(shared)
	diamond.Rule(f2).Run(shared)

	cases := []struct {
		name      string
		root      *Graph
		wantCycle bool
	}{
		{"acyclic", acyclic, false},
		{"selfLoop", selfLoop, true},
		{"abaLoop", ga, true},
		{"diamond", diamond, false},
	}
	for _, tc := range cases {
		err := detectRunCycle(tc.root)
		if !tc.wantCycle && err != nil {
			t.Errorf("%s: expected nil (legal DAG), got %v", tc.name, err)
			continue
		}
		if !tc.wantCycle {
			continue
		}
		if err == nil {
			t.Errorf("%s: expected cycle error, got nil", tc.name)
			continue
		}
		if !strings.Contains(err.Error(), "cycle") {
			t.Errorf("%s: error should mention cycle: %v", tc.name, err)
		}
	}
}
