//ff:func feature=policy type=rule control=sequence
//ff:what IsIPInList: backing(IPListBacking)으로 전달된 IP 목록에 클라이언트 IP가 있는지 판정
package policy

// IsIPInList checks if the client IP is in the list provided by backing.
// backing is *IPListBacking — purpose identifies the list, check performs the lookup.
func IsIPInList(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RequestContext)
	b := backing.(*IPListBacking)
	return b.Check(ctx.ClientIP), nil
}
