//ff:func feature=feature type=adapter control=sequence
//ff:what Inject: 활성 피처 목록을 request context에 주입
package feature

import (
	"context"
	"net/http"
)

type featuresKey struct{}

// Inject returns a net/http middleware that stores the list of enabled features
// in the request context under the featuresKey.
func Inject(f *Flags, ctxFn ContextFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uc := ctxFn(r)
			active, _ := f.List(uc)
			ctx := context.WithValue(r.Context(), featuresKey{}, active)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ActiveFeatures retrieves the list of enabled feature names from the request context.
func ActiveFeatures(r *http.Request) []string {
	v, _ := r.Context().Value(featuresKey{}).([]string)
	return v
}
