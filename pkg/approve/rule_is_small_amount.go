//ff:func feature=approve type=rule control=sequence
//ff:what IsSmallAmount: backing(*ThresholdBacking)으로 지정된 금액 임계값 이하인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsSmallAmount checks if the requested amount is at or below backing threshold.
func IsSmallAmount(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	amount, _ := ctx.Get("amount")
	b := backing.(*ThresholdBacking)
	return amount.(float64) <= b.Max, nil
}
