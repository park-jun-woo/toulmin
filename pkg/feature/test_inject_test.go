//ff:func feature=feature type=engine control=sequence
//ff:what TestInject — tests that Inject middleware stores active features in request context
package feature

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestInject(t *testing.T) {
	flags := NewFlags()
	g := toulmin.NewGraph("feature:dark-mode")
	g.Rule(IsBetaUser)
	flags.Register("dark-mode", g)

	ctxFn := func(r *http.Request) *UserContext {
		return &UserContext{Attributes: map[string]any{"beta": true}}
	}

	var gotFeatures []string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotFeatures = ActiveFeatures(r)
		w.WriteHeader(http.StatusOK)
	})

	handler := Inject(flags, ctxFn)(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if len(gotFeatures) != 1 || gotFeatures[0] != "dark-mode" {
		t.Fatalf("expected [dark-mode], got %v", gotFeatures)
	}
}
