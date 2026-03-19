package policy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

var testRoleFunc = func(u any) string { return u.(*testUser).Role }

func buildTestCtxFn(user any, ip string, headers map[string]string) ContextFunc {
	return func(r *http.Request) *RequestContext {
		return &RequestContext{
			User:     user,
			ClientIP: ip,
			Headers:  headers,
		}
	}
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok":true}`))
})

func TestGuard_AuthenticatedAdmin(t *testing.T) {
	g := toulmin.NewGraph("test:admin")
	g.Warrant(IsAuthenticated, nil, 1.0)
	g.Warrant(IsInRole, &RoleBacking{Role: "admin", RoleFunc: testRoleFunc}, 1.0)

	handler := Guard(g, buildTestCtxFn(&testUser{ID: "u1", Role: "admin"}, "10.0.0.1", nil))(okHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	handler.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

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

func TestGuard_IPBlocked(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", Check: func(ip string) bool { return ip == "1.2.3.4" }}

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

func TestGuardDebug_Forbidden(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", Check: func(ip string) bool { return ip == "1.2.3.4" }}

	g := toulmin.NewGraph("test:debug-deny")
	auth := g.Warrant(IsAuthenticated, nil, 1.0)
	blocked := g.Rebuttal(IsIPInList, blocklist, 1.0)
	g.Defeat(blocked, auth)

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
