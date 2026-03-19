//ff:func feature=approve type=rule control=sequence
//ff:what IsBudgetFrozen: 예산이 동결되었는지 판정
package approve

// IsBudgetFrozen checks if the budget is frozen.
func IsBudgetFrozen(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ApprovalContext)
	return ctx.Budget.Frozen, nil
}
