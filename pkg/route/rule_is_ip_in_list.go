//ff:func feature=route type=rule control=sequence
//ff:what IsIPInList: backing(func)으로 전달된 IP 목록에 클라이언트 IP가 있는지 판정
package route

// IsIPInList checks if the client IP is in the list provided by backing.
// backing is func(string) bool — same function for blocklist and whitelist,
// differentiated by backing value at graph declaration.
func IsIPInList(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	checkList := backing.(func(string) bool)
	return checkList(ctx.ClientIP), nil
}
