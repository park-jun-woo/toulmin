//ff:func feature=policy type=engine control=sequence
//ff:what TestGuardDebug_Branches — tests GuardDebug returns 500 on graph evaluation error (cyclic defeats) and 403 when no results
package policy

import (
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuardDebug_Branches(t *testing.T) {
	t.Run("eval error", func(t *testing.T) {
		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }

		g := toulmin.NewGraph("test:debug-cycle")
		ca := g.Rule(cycleA)
		cb := g.Counter(cycleB)
		cb.Attacks(ca)
		ca.Attacks(cb)

		handler := GuardDebug(g, buildTestCtxFn(&testUser{ID: "u1"}, "10.0.0.1", nil))(okHandler)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/debug", nil)
		handler.ServeHTTP(w, req)

		if w.Code != 500 {
			t.Errorf("expected 500 for evaluation error, got %d", w.Code)
		}
	})

	t.Run("no results", func(t *testing.T) {
		g := toulmin.NewGraph("test:debug-empty")

		handler := GuardDebug(g, buildTestCtxFn(&testUser{ID: "u1"}, "10.0.0.1", nil))(okHandler)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/debug", nil)
		handler.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Errorf("expected 403 for empty results, got %d", w.Code)
		}
	})
}
