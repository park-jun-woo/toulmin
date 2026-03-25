//ff:func feature=approve type=rule control=sequence
//ff:what IsBudgetFrozen: 예산이 동결되었는지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsBudgetFrozen checks if the budget is frozen.
func IsBudgetFrozen(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ApprovalContext)
	return ctx.Budget.Frozen, nil
}
