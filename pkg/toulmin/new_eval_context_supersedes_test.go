//ff:func feature=engine type=engine control=sequence
//ff:what TestDetectSupersedeCycle — supersede cycles are rejected independently of the Attacks graph
package toulmin

import (
	"strings"
	"testing"
)

func TestDetectSupersedeCycle(t *testing.T) {
	// Mutually-excluding supersede specs (X excludes Y / Y excludes X) are rejected.
	cyc := []defeatEdge{
		{from: "X", to: "Y"},
		{from: "Y", to: "X"},
	}
	err := detectSupersedeCycle(cyc)
	if err == nil {
		t.Fatalf("expected supersede cycle error")
	}
	if !strings.Contains(err.Error(), "toulmin: supersede cycle detected at") {
		t.Fatalf("expected wrapped supersede cycle message, got %q", err.Error())
	}

	// An acyclic supersede graph passes.
	if err := detectSupersedeCycle([]defeatEdge{{from: "X", to: "Y"}}); err != nil {
		t.Fatalf("acyclic supersede graph should pass, got %v", err)
	}

	// A supersede cycle must not leak into the Attacks cycle check: a graph with a
	// clean Attacks graph but a cyclic supersede graph fails on supersede only, and
	// a graph with a clean supersede graph but cyclic Attacks fails on Attacks only.
	rules := []RuleMeta{
		{Name: "X", Qualifier: 1.0, Strength: Defeasible, Fn: WarrantA},
		{Name: "Y", Qualifier: 1.0, Strength: Defeasible, Fn: RebuttalB},
	}
	// clean Attacks (nil), cyclic supersede → supersede error
	_, err = newEvalContext(rules, nil, cyc, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "supersede cycle detected") {
		t.Fatalf("expected supersede cycle error, got %v", err)
	}
	// cyclic Attacks, clean supersede (nil) → plain Attacks cycle error (not wrapped)
	attackCyc := []defeatEdge{{from: "X", to: "Y"}, {from: "Y", to: "X"}}
	_, err = newEvalContext(rules, attackCyc, nil, map[string]string{})
	if err == nil || strings.Contains(err.Error(), "supersede") {
		t.Fatalf("expected plain Attacks cycle error, got %v", err)
	}
}
