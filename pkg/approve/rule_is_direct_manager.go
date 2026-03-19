//ff:func feature=approve type=rule control=sequence
//ff:what IsDirectManager: 결재자가 요청자의 직속 상위자인지 판정
package approve

// IsDirectManager checks if the approver is the requester's direct manager.
// backing is *ApproverBacking — IDFunc extracts approver ID.
func IsDirectManager(claim any, ground any, backing any) (bool, any) {
	req := claim.(*ApprovalRequest)
	ctx := ground.(*ApprovalContext)
	ab := backing.(*ApproverBacking)
	return ctx.OrgTree.IsDirectManager(ab.IDFunc(ctx.Approver), req.RequesterID), nil
}
