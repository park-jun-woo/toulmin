//ff:func feature=approve type=rule control=sequence
//ff:what IsDirectManager: 결재자가 요청자의 직속 상위자인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsDirectManager checks if the approver is the requester's direct manager.
func IsDirectManager(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	orgTree, _ := ctx.Get("orgTree")
	approverID, _ := ctx.Get("approverID")
	requesterID, _ := ctx.Get("requesterID")
	return orgTree.(OrgTree).IsDirectManager(approverID.(string), requesterID.(string)), nil
}
