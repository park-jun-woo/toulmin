//ff:func feature=policy type=adapter control=sequence
//ff:what Guard: toulmin graph를 net/http 정책 미들웨어로 변환
package policy

import (
	"net/http"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Guard returns a net/http middleware that evaluates the policy graph.
// Uses Evaluate (lightweight). verdict <= 0 is denied with 403.
func Guard(g *toulmin.Graph, ctxFn ContextFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxFn(r)
			results, err := g.Evaluate(nil, ctx)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"policy evaluation failed"}`))
				return
			}
			if len(results) == 0 || results[0].Verdict <= 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"forbidden"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
