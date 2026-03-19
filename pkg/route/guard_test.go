package route

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
	return func(c *gin.Context) *RouteContext {
		return &RouteContext{
			User:     user,
			ClientIP: ip,
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
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
	blocklist := func(ip string) bool { return ip == "1.2.3.4" }

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
	blocklist := func(ip string) bool { return ip == "1.2.3.4" }
	whitelist := func(ip string) bool { return ip == "1.2.3.4" }

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

func TestGuard_RateLimited_InternalServiceDefeat(t *testing.T) {
	limiter := &mockLimiter{limited: map[string]bool{"10.0.0.1": true}}

	g := toulmin.NewGraph("test:ratelimit")
	auth := g.Warrant(IsAuthenticated, nil, 1.0)
	limited := g.Rebuttal(IsRateLimited, limiter, 1.0)
	internal := g.Defeater(IsInternalService, nil, 1.0)
	g.Defeat(limited, auth)
	g.Defeat(internal, limited)

	headers := map[string]string{"X-Internal-Token": "secret"}
	r := gin.New()
	r.GET("/api", Guard(g, buildTestCtx(&User{ID: "u1"}, "10.0.0.1", headers)), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200 (internal service defeats rate limit), got %d", w.Code)
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
	verdict := w.Header().Get("X-Route-Verdict")
	if verdict == "" {
		t.Error("expected X-Route-Verdict header")
	}
	trace := w.Header().Get("X-Route-Trace")
	if trace == "" {
		t.Error("expected X-Route-Trace header")
	}
}

func TestGuardDebug_Forbidden_WithRebuttal(t *testing.T) {
	blocklist := func(ip string) bool { return ip == "1.2.3.4" }

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
	verdict := w.Header().Get("X-Route-Verdict")
	if verdict == "" {
		t.Error("expected X-Route-Verdict header even on deny")
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "forbidden" {
		t.Errorf("expected forbidden error, got %v", body)
	}
}
