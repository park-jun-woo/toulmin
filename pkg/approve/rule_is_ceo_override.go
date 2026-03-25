//ff:func feature=approve type=rule control=sequence
//ff:what IsCEOOverride: CEO 직권 승인 여부 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsCEOOverride checks if the approver is the CEO.
func IsCEOOverride(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ApprovalContext)
	return ctx.ApproverRole == "ceo", nil
}
