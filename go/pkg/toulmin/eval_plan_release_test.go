//ff:func feature=engine type=engine control=sequence
//ff:what TestEvalPlanRelease — release: every per-evaluation map, the trace and the error of the state are cleared, and the shared maps are left as they are
package toulmin

import (
	"errors"
	"testing"
)

func TestEvalPlanRelease(t *testing.T) {
	g := NewGraph("rel")
	g.Rule(func(ctx Context, specs Specs) (float64, any) { return 1.0, "ev" })
	p := g.plan()
	ec, _ := p.context()
	ec.ran["a"], ec.active["a"], ec.ruleVerdict["a"], ec.applicability["a"] = true, true, 1, 1
	ec.evidence["a"], ec.verdictCache["a"] = "ev", 1
	ec.trace, ec.err = []TraceEntry{{}}, errors.New("x")
	p.release(ec)
	if len(ec.ran)+len(ec.active)+len(ec.ruleVerdict)+len(ec.applicability)+len(ec.evidence)+len(ec.verdictCache) != 0 || ec.trace != nil || ec.err != nil {
		t.Fatalf("expected a cleared state, got %+v", ec)
	}
	if len(ec.fnMap) != 1 || ec.home != p {
		t.Fatal("the shared maps stay")
	}
}
