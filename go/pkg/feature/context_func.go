//ff:type feature=feature type=adapter
//ff:what ContextFunc: http.Request → UserContext 변환 함수 타입
package feature

import "net/http"

// ContextFunc converts an http.Request into a UserContext.
type ContextFunc func(r *http.Request) *UserContext
