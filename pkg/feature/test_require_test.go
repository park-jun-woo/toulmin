//ff:func feature=feature type=engine control=sequence
//ff:what TestRequire — tests Require middleware gating on feature enablement
package feature

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestRequire(t *testing.T) {
	t.Run("Enabled", func(t *testing.T) {
		flags := NewFlags()
		g := toulmin.NewGraph("feature:dark-mode")
		g.Rule(IsBetaUser)
		flags.Register("dark-mode", g)

		ctxFn := func(r *http.Request) *UserContext {
			return &UserContext{Attributes: map[string]any{"beta": true}}
		}

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		handler := Require(flags, "dark-mode", ctxFn)(next)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if !called {
			t.Fatal("expected next handler to be called")
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		flags := NewFlags()
		g := toulmin.NewGraph("feature:dark-mode")
		g.Rule(IsBetaUser)
		flags.Register("dark-mode", g)

		ctxFn := func(r *http.Request) *UserContext {
			return &UserContext{Attributes: map[string]any{"beta": false}}
		}

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		handler := Require(flags, "dark-mode", ctxFn)(next)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if called {
			t.Fatal("expected next handler not to be called")
		}
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("Error", func(t *testing.T) {
		flags := NewFlags()

		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		g := toulmin.NewGraph("feature:cycle")
		a := g.Rule(cycleA)
		b := g.Counter(cycleB)
		b.Attacks(a)
		a.Attacks(b)
		flags.Register("cycle", g)

		ctxFn := func(r *http.Request) *UserContext {
			return &UserContext{}
		}

		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		handler := Require(flags, "cycle", ctxFn)(next)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if called {
			t.Fatal("expected next handler not to be called")
		}
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", rec.Code)
		}
	})
}
