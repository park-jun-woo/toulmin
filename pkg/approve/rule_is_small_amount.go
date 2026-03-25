//ff:func feature=approve type=rule control=sequence
//ff:what IsSmallAmount: spec(*ThresholdSpec)으로 지정된 금액 임계값 이하인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsSmallAmount checks if the requested amount is at or below spec threshold.
func IsSmallAmount(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	amount, _ := ctx.Get("amount")
	b := specs[0].(*ThresholdSpec)
	return amount.(float64) <= b.Max, nil
}
