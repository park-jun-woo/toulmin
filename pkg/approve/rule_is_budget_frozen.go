//ff:func feature=approve type=rule control=sequence
//ff:what IsBudgetFrozen: 예산이 동결되었는지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBudgetFrozen checks if the budget is frozen.
func IsBudgetFrozen(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	budget, _ := ctx.Get("budget")
	b, ok := budget.(*Budget)
	if !ok {
		return false, nil
	}
	return b.Frozen, nil
}
