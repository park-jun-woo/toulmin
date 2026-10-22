//ff:type feature=moderate type=adapter
//ff:what ExtractFunc: http.Request에서 Content와 ContentContext를 추출하는 함수 타입
package moderate

import "net/http"

// ExtractFunc extracts Content and ContentContext from an http.Request.
type ExtractFunc func(r *http.Request) (*Content, *ContentContext)
