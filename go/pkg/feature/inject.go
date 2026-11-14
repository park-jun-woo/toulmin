//ff:func feature=feature type=adapter control=sequence
//ff:what Inject — net/http middleware that stores active features in request context
package feature

import (
	"context"
	"net/http"
)

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
