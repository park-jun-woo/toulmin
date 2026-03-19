//ff:func feature=feature type=adapter control=sequence
//ff:what Require: 특정 피처가 활성화된 사용자만 접근 허용
package feature

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Require returns a gin.HandlerFunc that allows access only if the feature is enabled.
func Require(f *Flags, name string, ctxBuilder FeatureContextFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxBuilder(c)
		enabled, err := f.IsEnabled(name, ctx)
		if err != nil || !enabled {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}
