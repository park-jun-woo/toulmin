//ff:func feature=policy type=engine control=sequence
//ff:what buildTestCtxFn — builds test context function for guard tests
package policy

import "net/http"

func buildTestCtxFn(user any, ip string, headers map[string]string) ContextFunc {
	role := ""
	if u, ok := user.(*testUser); ok {
		role = u.Role
	}
	return func(r *http.Request) *RequestContext {
		return &RequestContext{
			User:     user,
			ClientIP: ip,
			Headers:  headers,
			Role:     role,
		}
	}
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok":true}`))
})
