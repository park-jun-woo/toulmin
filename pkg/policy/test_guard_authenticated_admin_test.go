//ff:func feature=policy type=engine control=sequence
//ff:what TestGuard_AuthenticatedAdmin — tests Guard allows authenticated admin
package policy

import (
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuard_AuthenticatedAdmin(t *testing.T) {
	g := toulmin.NewGraph("test:admin")
	g.Rule(IsAuthenticated)
	g.Rule(IsInRole).Backing(&RoleBacking{Role: "admin"})

	handler := Guard(g, buildTestCtxFn(&testUser{ID: "u1", Role: "admin"}, "10.0.0.1", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
