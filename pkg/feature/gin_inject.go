//ff:func feature=feature type=adapter control=sequence
//ff:what Inject: 활성 피처 목록을 gin.Context에 주입
package feature

import "github.com/gin-gonic/gin"

// Inject returns a gin.HandlerFunc that stores the list of enabled features in gin.Context.
func Inject(f *Flags, ctxBuilder FeatureContextFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxBuilder(c)
		active, _ := f.List(ctx)
		c.Set("features", active)
		c.Next()
	}
}
