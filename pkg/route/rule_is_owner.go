//ff:func feature=route type=rule control=sequence
//ff:what IsOwner: backing(func)으로 추출한 소유자 ID와 사용자 ID 비교
package route

// IsOwner checks if the user owns the resource.
// backing is func(*RouteContext) string that extracts the resource owner ID.
func IsOwner(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*RouteContext)
	if ctx.User == nil {
		return false, nil
	}
	ownerIDFunc := backing.(func(*RouteContext) string)
	return ctx.User.ID == ownerIDFunc(ctx), nil
}
