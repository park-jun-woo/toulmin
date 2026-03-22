//ff:func feature=policy type=engine control=sequence
//ff:what buildTestCtxFn — builds test context function for guard tests
package policy

import "net/http"

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
