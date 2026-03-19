//ff:func feature=route type=adapter control=sequence
//ff:what Guard: toulmin graph를 Gin 미들웨어로 변환하는 어댑터
package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// ContextBuilderFunc converts a gin.Context into a RouteContext.
type ContextBuilderFunc func(*gin.Context) *RouteContext

// Guard returns a gin.HandlerFunc that evaluates the given graph.
// Uses Evaluate (lightweight, no trace). For debug trace, use GuardDebug.
// claim is nil — route matching is already done by Gin.
// verdict <= 0 is denied (security context: undecided is deny).
func Guard(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxBuilder(c)
		results, err := g.Evaluate(nil, ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "route evaluation failed"})
			return
		}
		if len(results) == 0 || results[0].Verdict <= 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
