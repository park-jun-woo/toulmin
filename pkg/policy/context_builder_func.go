//ff:type feature=policy type=adapter
//ff:what ContextBuilderFunc: gin.Context → RequestContext 변환 함수 타입
package policy

import "github.com/gin-gonic/gin"

// ContextBuilderFunc converts a gin.Context into a RequestContext.
type ContextBuilderFunc func(*gin.Context) *RequestContext
