//ff:func feature=policy type=adapter control=sequence
//ff:what GuardDebug: 판정 근거를 응답 헤더로 노출하는 디버그 미들웨어
package policy

import (
	"encoding/json"
	"net/http"

	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// GuardDebug returns a net/http middleware that evaluates the policy graph
// and exposes verdict and trace in response headers and body.
func GuardDebug(g *toulmin.Graph, ctxFn ContextFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ctxFn(r)
			results, err := g.Evaluate(nil, ctx, toulmin.EvalOption{Trace: true})
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"policy evaluation failed"}`))
				return
			}
			if len(results) == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"forbidden"}`))
				return
			}
			w.Header().Set("X-Policy-Verdict", formatVerdict(results[0].Verdict))
			w.Header().Set("X-Policy-Trace", formatTrace(results[0].Trace))
			if results[0].Verdict <= 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				body, _ := json.Marshal(map[string]string{
					"error": "forbidden",
					"trace": formatTrace(results[0].Trace),
				})
				w.Write(body)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
