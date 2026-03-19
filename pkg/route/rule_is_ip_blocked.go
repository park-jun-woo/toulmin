//ff:func feature=route type=rule control=sequence
//ff:what IsIPBlocked: backing(func)으로 전달된 차단 목록에 IP가 있는지 판정
package route

// IsIPBlocked checks if the client IP is in the blocklist.
// backing is func(string) bool that checks the blocklist.
func IsIPBlocked(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	blocklist := backing.(func(string) bool)
	return blocklist(ctx.ClientIP), nil
}
