package policy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

var testRoleFunc = func(u any) string { return u.(*testUser).Role }

func buildTestCtx(user any, ip string, headers map[string]string) ContextBuilderFunc {
	return func(c *gin.Context) *RequestContext {
		return &RequestContext{
			User:     user,
			ClientIP: ip,
			Headers:  headers,
		}
	}
}

func TestGuard_AuthenticatedAdmin(t *testing.T) {
	g := toulmin.NewGraph("test:admin")
	g.Warrant(IsAuthenticated, nil, 1.0)
	g.Warrant(IsInRole, &RoleBacking{Role: "admin", RoleFunc: testRoleFunc}, 1.0)

	r := gin.New()
	r.GET("/admin", Guard(g, buildTestCtx(&testUser{ID: "u1", Role: "admin"}, "10.0.0.1", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGuard_Unauthenticated(t *testing.T) {
	g := toulmin.NewGraph("test:auth")
	g.Warrant(IsAuthenticated, nil, 1.0)

	r := gin.New()
	r.GET("/protected", Guard(g, buildTestCtx(nil, "10.0.0.1", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	r.ServeHTTP(w, req)

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

	r := gin.New()
	r.GET("/api", Guard(g, buildTestCtx(&testUser{ID: "u1"}, "1.2.3.4", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403 for blocked IP, got %d", w.Code)
	}
}

func TestGuardDebug_Headers(t *testing.T) {
	g := toulmin.NewGraph("test:debug")
	g.Warrant(IsAuthenticated, nil, 1.0)

	r := gin.New()
	r.GET("/debug", GuardDebug(g, buildTestCtx(&testUser{ID: "u1"}, "10.0.0.1", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/debug", nil)
	r.ServeHTTP(w, req)

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

	r := gin.New()
	r.GET("/debug", GuardDebug(g, buildTestCtx(&testUser{ID: "u1"}, "1.2.3.4", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/debug", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}

	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "forbidden" {
		t.Errorf("expected forbidden error, got %v", body)
	}
}
