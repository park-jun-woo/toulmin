//ff:func feature=approve type=rule control=sequence
//ff:what IsDirectManager: 결재자가 요청자의 직속 상위자인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsDirectManager checks if the approver is the requester's direct manager.
func IsDirectManager(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	orgTree, _ := ctx.Get("orgTree")
	approverID, _ := ctx.Get("approverID")
	requesterID, _ := ctx.Get("requesterID")
	ot, ok := orgTree.(OrgTree)
	if !ok {
		return false, nil
	}
	aid, ok := approverID.(string)
	if !ok {
		return false, nil
	}
	rid, ok := requesterID.(string)
	if !ok {
		return false, nil
	}
	return ot.IsDirectManager(aid, rid), nil
}
