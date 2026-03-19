//ff:func feature=route type=adapter control=sequence
//ff:what GuardDebug: 판정 근거를 응답 헤더로 노출하는 디버그 미들웨어
package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// GuardDebug returns a gin.HandlerFunc that evaluates the given graph
// and exposes verdict and trace in response headers.
func GuardDebug(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxBuilder(c)
		results, err := g.EvaluateTrace(nil, ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "route evaluation failed"})
			return
		}
		if len(results) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Header("X-Route-Verdict", fmt.Sprintf("%.2f", results[0].Verdict))
		c.Header("X-Route-Trace", formatTrace(results[0].Trace))
		if results[0].Verdict <= 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func formatTrace(traces []toulmin.TraceEntry) string {
	parts := make([]string, len(traces))
	for i, t := range traces {
		parts[i] = fmt.Sprintf("%s=%v", t.Name, t.Activated)
	}
	return strings.Join(parts, ", ")
}
