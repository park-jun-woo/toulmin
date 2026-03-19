//ff:func feature=policy type=adapter control=sequence
//ff:what GuardDebug: 판정 근거를 응답 헤더로 노출하는 디버그 미들웨어
package policy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// GuardDebug returns a gin.HandlerFunc that evaluates the policy graph
// and exposes verdict and trace in response headers and body.
func GuardDebug(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxBuilder(c)
		results, err := g.EvaluateTrace(nil, ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "policy evaluation failed"})
			return
		}
		if len(results) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Header("X-Policy-Verdict", formatVerdict(results[0].Verdict))
		c.Header("X-Policy-Trace", formatTrace(results[0].Trace))
		if results[0].Verdict <= 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
				"trace": formatTrace(results[0].Trace),
			})
			return
		}
		c.Next()
	}
}
