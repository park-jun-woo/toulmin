//ff:func feature=approve type=rule control=sequence
//ff:what IsAboveLevel: backing(ApproverBacking)으로 지정된 최소 레벨 이상인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsAboveLevel checks if the approver's level is at or above backing.Level.
func IsAboveLevel(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ApprovalContext)
	ab := backing.(*ApproverBacking)
	return ctx.ApproverLevel >= ab.Level, nil
}
