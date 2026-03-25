//ff:func feature=approve type=rule control=sequence
//ff:what IsDirectManager: 결재자가 요청자의 직속 상위자인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsDirectManager checks if the approver is the requester's direct manager.
func IsDirectManager(claim any, ground any, backing toulmin.Backing) (bool, any) {
	req := claim.(*ApprovalRequest)
	ctx := ground.(*ApprovalContext)
	return ctx.OrgTree.IsDirectManager(ctx.ApproverID, req.RequesterID), nil
}
