//ff:func feature=approve type=rule control=sequence
//ff:what IsCEOOverride: CEO 직권 승인 여부 판정
package approve

// IsCEOOverride checks if the approver is the CEO.
func IsCEOOverride(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ApprovalContext)
	return ctx.Approver.Role == "ceo", nil
}
