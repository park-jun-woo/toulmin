//ff:func feature=route type=rule control=sequence
//ff:what IsInternalService: 내부 서비스 요청인지 판정
package route

// IsInternalService returns true if the request has an internal service token header.
func IsInternalService(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	token := ctx.Headers["X-Internal-Token"]
	return token != "", nil
}
