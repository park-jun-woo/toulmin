//ff:func feature=approve type=rule control=sequence
//ff:what IsUnderBudget: 요청 금액이 잔여 예산 이하인지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsUnderBudget checks if the requested amount is within remaining budget.
func IsUnderBudget(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	amount, _ := ctx.Get("amount")
	budget, _ := ctx.Get("budget")
	a, ok := amount.(float64)
	if !ok {
		return false, nil
	}
	b, ok := budget.(*Budget)
	if !ok {
		return false, nil
	}
	return a <= b.Remaining, nil
}
