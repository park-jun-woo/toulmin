//ff:func feature=route type=rule control=sequence
//ff:what IsWhitelisted: backing(func)으로 전달된 화이트리스트에 IP가 있는지 판정
package route

// IsWhitelisted checks if the client IP is in the whitelist.
// backing is func(string) bool that checks the whitelist.
func IsWhitelisted(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	whitelist := backing.(func(string) bool)
	return whitelist(ctx.ClientIP), nil
}
