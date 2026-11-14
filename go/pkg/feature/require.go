//ff:func feature=feature type=adapter control=sequence
//ff:what Require: 특정 피처가 활성화된 사용자만 접근 허용
package feature

import "net/http"

// Require returns a net/http middleware that allows access only if the feature is enabled.
// Returns 404 if the feature is disabled or an error occurs.
func Require(f *Flags, name string, ctxFn ContextFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxFn(r)
			enabled, err := f.IsEnabled(name, ctx)
			if err != nil || !enabled {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
