//ff:func feature=approve type=rule control=sequence
//ff:what IsAboveLevel: backing(int)으로 지정된 최소 레벨 이상인지 판정
package approve

// IsAboveLevel checks if the approver's level is at or above backing (int).
func IsAboveLevel(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ApprovalContext)
	minLevel := backing.(int)
	return ctx.Approver.Level >= minLevel, nil
}
