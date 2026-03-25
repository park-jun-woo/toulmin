//ff:func feature=policy type=engine control=sequence
//ff:what TestGuardDebug_Forbidden — tests GuardDebug returns forbidden with trace
package policy

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func TestGuardDebug_Forbidden(t *testing.T) {
	blocklist := &IPListSpec{Purpose: "blocklist", List: []string{"1.2.3.4"}}

	g := toulmin.NewGraph("test:debug-deny")
	auth := g.Rule(IsAuthenticated)
	blocked := g.Counter(IsIPInList).With(blocklist)
	blocked.Attacks(auth)

	handler := GuardDebug(g, buildTestCtxFn(&testUser{ID: "u1"}, "1.2.3.4", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/debug", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}

	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "forbidden" {
		t.Errorf("expected forbidden error, got %v", body)
	}
}
