//ff:func feature=policy type=adapter control=sequence
//ff:what Guard: toulmin graph를 Gin 정책 미들웨어로 변환
package policy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Guard returns a gin.HandlerFunc that evaluates the policy graph.
// Uses Evaluate (lightweight). verdict <= 0 is denied.
func Guard(g *toulmin.Graph, ctxBuilder ContextBuilderFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxBuilder(c)
		results, err := g.Evaluate(nil, ctx)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "policy evaluation failed"})
			return
		}
		if len(results) == 0 || results[0].Verdict <= 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
