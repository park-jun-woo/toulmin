//ff:type feature=policy type=adapter
//ff:what ContextFunc: http.Request → RequestContext 변환 함수 타입
package policy

import "net/http"

// ContextFunc converts an http.Request into a RequestContext.
type ContextFunc func(r *http.Request) *RequestContext
