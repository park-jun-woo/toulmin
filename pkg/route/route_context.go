//ff:type feature=route type=model
//ff:what RouteContext: 라우트 판정에 필요한 요청 컨텍스트
package route

// RouteContext holds request context for route evaluation.
type RouteContext struct {
	User     *User
	ClientIP string
	Method   string
	Path     string
	Headers  map[string]string
	Metadata map[string]any
}
