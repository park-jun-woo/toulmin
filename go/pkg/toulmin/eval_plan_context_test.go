//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalPlanContext — context: the plan's error for a cyclic graph; otherwise a state that shares the plan's maps, knows its plan and has empty per-evaluation maps, new when none was released and the released one after a release
package toulmin

import (
	"errors"
	"testing"
)

func TestEvalPlanContext(t *testing.T) {
	bad := &evalPlan{err: errors.New("cycle")}
	if ec, err := bad.context(); ec != nil || err == nil {
		t.Fatalf("expected the plan's error, got %v, %v", ec, err)
	}
	g := NewGraph("ctx")
	g.Rule(func(ctx Context, specs Specs) (float64, any) { return 1.0, nil })
	p := g.plan()
	a, err := p.context()
	if err != nil || a == p.static || a.home != p || len(a.fnMap) != 1 || a.ran == nil || len(a.ran) != 0 || a.verdictCache == nil {
		t.Fatalf("expected a fresh state on the plan's maps, got %+v (%v)", a, err)
	}
	if b, _ := p.context(); b == a {
		t.Fatal("a state in use is not handed out again")
	}
	reused := false
	for i := 0; i < 100 && !reused; i++ {
		a.ran["x"] = true
		p.release(a)
		b, _ := p.context()
		reused = b == a && len(b.ran) == 0
		a = b
	}
	if !reused {
		t.Fatal("a released state is handed out again, cleared")
	}
}
