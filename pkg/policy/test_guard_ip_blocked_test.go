//ff:func feature=policy type=engine control=sequence
//ff:what TestGuard_IPBlocked — tests Guard blocks request from blocked IP
package policy

import (
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuard_IPBlocked(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", List: []string{"1.2.3.4"}}

	g := toulmin.NewGraph("test:ip")
	auth := g.Warrant(IsAuthenticated, nil, 1.0)
	blocked := g.Rebuttal(IsIPInList, blocklist, 1.0)
	g.Defeat(blocked, auth)

	handler := Guard(g, buildTestCtxFn(&testUser{ID: "u1"}, "1.2.3.4", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403 for blocked IP, got %d", w.Code)
	}
}
