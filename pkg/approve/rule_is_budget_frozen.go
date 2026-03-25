//ff:func feature=approve type=rule control=sequence
//ff:what IsBudgetFrozen: 예산이 동결되었는지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBudgetFrozen checks if the budget is frozen.
func IsBudgetFrozen(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	budget, _ := ctx.Get("budget")
	return budget.(*Budget).Frozen, nil
}
