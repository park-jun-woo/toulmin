//ff:type feature=moderate type=adapter
//ff:what ContentExtractor: gin.Context에서 Content와 ContentContext를 추출하는 함수 타입
package moderate

import "github.com/gin-gonic/gin"

// ContentExtractor extracts Content and ContentContext from a gin.Context.
type ContentExtractor func(*gin.Context) (*Content, *ContentContext)
