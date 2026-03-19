//ff:type feature=feature type=adapter
//ff:what FeatureContextFunc: gin.Context → UserContext 변환 함수 타입
package feature

import "github.com/gin-gonic/gin"

// FeatureContextFunc converts a gin.Context into a UserContext.
type FeatureContextFunc func(*gin.Context) *UserContext
