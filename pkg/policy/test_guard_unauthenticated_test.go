//ff:func feature=policy type=engine control=sequence
//ff:what TestGuard_Unauthenticated — tests Guard blocks unauthenticated request
package policy

import (
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuard_Unauthenticated(t *testing.T) {
	g := toulmin.NewGraph("test:auth")
	g.Warrant(IsAuthenticated, nil, 1.0)

	handler := Guard(g, buildTestCtxFn(nil, "10.0.0.1", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
