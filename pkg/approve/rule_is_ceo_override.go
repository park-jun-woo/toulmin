//ff:func feature=approve type=rule control=sequence
//ff:what IsCEOOverride: CEO 직권 승인 여부 판정
package approve

// IsCEOOverride checks if the approver is the CEO.
// backing is *ApproverBacking — RoleFunc extracts approver role.
func IsCEOOverride(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ApprovalContext)
	ab := backing.(*ApproverBacking)
	return ab.RoleFunc(ctx.Approver) == "ceo", nil
}
