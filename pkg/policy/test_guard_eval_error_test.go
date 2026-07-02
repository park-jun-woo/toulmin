//ff:func feature=policy type=engine control=sequence
//ff:what TestGuard_EvalError — tests Guard returns 500 when graph evaluation errors (cyclic defeats)
package policy

import (
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuard_EvalError(t *testing.T) {
	cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
	cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }

	g := toulmin.NewGraph("test:cycle")
	ca := g.Rule(cycleA)
	cb := g.Counter(cycleB)
	cb.Attacks(ca)
	ca.Attacks(cb)

	handler := Guard(g, buildTestCtxFn(&testUser{ID: "u1"}, "10.0.0.1", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected 500 for evaluation error, got %d", w.Code)
	}
}
