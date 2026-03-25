//ff:func feature=approve type=rule control=sequence
//ff:what IsSmallAmount: backing(*ThresholdBacking)으로 지정된 금액 임계값 이하인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsSmallAmount checks if the requested amount is at or below backing threshold.
func IsSmallAmount(claim any, ground any, backing toulmin.Backing) (bool, any) {
	req := claim.(*ApprovalRequest)
	b := backing.(*ThresholdBacking)
	return req.Amount <= b.Max, nil
}
