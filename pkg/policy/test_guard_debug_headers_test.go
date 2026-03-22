//ff:func feature=policy type=engine control=sequence
//ff:what TestGuardDebug_Headers — tests GuardDebug adds verdict headers
package policy

import (
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuardDebug_Headers(t *testing.T) {
	g := toulmin.NewGraph("test:debug")
	g.Warrant(IsAuthenticated, nil, 1.0)

	handler := GuardDebug(g, buildTestCtxFn(&testUser{ID: "u1"}, "10.0.0.1", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/debug", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-Policy-Verdict") == "" {
		t.Error("expected X-Policy-Verdict header")
	}
}
