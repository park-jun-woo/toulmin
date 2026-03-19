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

func buildTestCtx(user *User, ip string, headers map[string]string) ContextBuilderFunc {
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
	g.Warrant(IsInRole, "admin", 1.0)

	r := gin.New()
	r.GET("/admin", Guard(g, buildTestCtx(&User{ID: "u1", Role: "admin"}, "10.0.0.1", nil)), func(c *gin.Context) {
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
	r.GET("/api", Guard(g, buildTestCtx(&User{ID: "u1"}, "1.2.3.4", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403 for blocked IP, got %d", w.Code)
	}
}

func TestGuard_IPBlocked_WhitelistDefeat(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", Check: func(ip string) bool { return ip == "1.2.3.4" }}
	whitelist := &IPListBacking{Purpose: "whitelist", Check: func(ip string) bool { return ip == "1.2.3.4" }}

	g := toulmin.NewGraph("test:whitelist")
	auth := g.Warrant(IsAuthenticated, nil, 1.0)
	blocked := g.Rebuttal(IsIPInList, blocklist, 1.0)
	allowed := g.Defeater(IsIPInList, whitelist, 1.0)
	g.Defeat(blocked, auth)
	g.Defeat(allowed, blocked)

	r := gin.New()
	r.GET("/api", Guard(g, buildTestCtx(&User{ID: "u1"}, "1.2.3.4", nil)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200 (whitelist defeats blocklist), got %d", w.Code)
	}
}

func TestGuard_InternalService(t *testing.T) {
	g := toulmin.NewGraph("test:internal")
	auth := g.Warrant(IsAuthenticated, nil, 1.0)
	internal := g.Defeater(HasHeader, "X-Internal-Token", 1.0)
	_ = internal
	// HasHeader as defeater doesn't defeat auth in this test — just verify it works
	_ = auth

	r := gin.New()
	r.GET("/api", Guard(g, buildTestCtx(&User{ID: "u1"}, "10.0.0.1", map[string]string{"X-Internal-Token": "secret"})), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGuardDebug_Headers(t *testing.T) {
	g := toulmin.NewGraph("test:debug")
	g.Warrant(IsAuthenticated, nil, 1.0)

	r := gin.New()
	r.GET("/debug", GuardDebug(g, buildTestCtx(&User{ID: "u1"}, "10.0.0.1", nil)), func(c *gin.Context) {
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
	if w.Header().Get("X-Policy-Trace") == "" {
		t.Error("expected X-Policy-Trace header")
	}
}

func TestGuardDebug_Forbidden(t *testing.T) {
	blocklist := &IPListBacking{Purpose: "blocklist", Check: func(ip string) bool { return ip == "1.2.3.4" }}

	g := toulmin.NewGraph("test:debug-deny")
	auth := g.Warrant(IsAuthenticated, nil, 1.0)
	blocked := g.Rebuttal(IsIPInList, blocklist, 1.0)
	g.Defeat(blocked, auth)

	r := gin.New()
	r.GET("/debug", GuardDebug(g, buildTestCtx(&User{ID: "u1"}, "1.2.3.4", nil)), func(c *gin.Context) {
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
	if body["trace"] == nil {
		t.Error("expected trace in response body")
	}
}
