//ff:func feature=moderate type=engine control=sequence
//ff:what TestGuard — tests Guard middleware branches for error, block, flag, and allow
package moderate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuard(t *testing.T) {
	run := func(mod *Moderator, extractFn ExtractFunc) (*httptest.ResponseRecorder, bool) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})
		handler := Guard(mod, extractFn)(next)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec, called
	}

	t.Run("Error", func(t *testing.T) {
		cycleA := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		cycleB := func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) { return true, nil }
		g := toulmin.NewGraph("test:cycle")
		a := g.Rule(cycleA)
		b := g.Counter(cycleB)
		b.Attacks(a)
		a.Attacks(b)
		mod := NewModerator(g)

		extractFn := func(r *http.Request) (*Content, *ContentContext) {
			return &Content{Body: "hi"}, &ContentContext{Author: &Author{}, Channel: &Channel{}}
		}

		rec, called := run(mod, extractFn)
		if called {
			t.Fatal("expected next handler not to be called")
		}
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", rec.Code)
		}
	})

	t.Run("Block", func(t *testing.T) {
		cb := &ClassifierSpec{Classifier: &mockClassifier{scores: map[string]float64{"hate_speech": 0.95}}}
		g := toulmin.NewGraph("test:block")
		verified := g.Rule(IsVerifiedUser)
		hate := g.Counter(ContainsHateSpeech).With(cb)
		hate.Attacks(verified)
		mod := NewModerator(g)

		extractFn := func(r *http.Request) (*Content, *ContentContext) {
			return &Content{Body: "hate content"}, &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "general"}}
		}

		rec, called := run(mod, extractFn)
		if called {
			t.Fatal("expected next handler not to be called")
		}
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", rec.Code)
		}
	})

	t.Run("Flag", func(t *testing.T) {
		g := toulmin.NewGraph("test:flag")
		g.Rule(IsVerifiedUser).Qualifier(0.6)
		mod := NewModerator(g)

		extractFn := func(r *http.Request) (*Content, *ContentContext) {
			return &Content{Body: "hi"}, &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "general"}}
		}

		rec, called := run(mod, extractFn)
		if called {
			t.Fatal("expected next handler not to be called")
		}
		if rec.Code != http.StatusAccepted {
			t.Fatalf("expected status 202, got %d", rec.Code)
		}
	})

	t.Run("Allow", func(t *testing.T) {
		g := toulmin.NewGraph("test:allow")
		g.Rule(IsVerifiedUser)
		mod := NewModerator(g)

		extractFn := func(r *http.Request) (*Content, *ContentContext) {
			return &Content{Body: "hi"}, &ContentContext{Author: &Author{Verified: true}, Channel: &Channel{Type: "general"}}
		}

		rec, called := run(mod, extractFn)
		if !called {
			t.Fatal("expected next handler to be called")
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}
	})
}
