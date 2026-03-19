//ff:func feature=approve type=rule control=sequence
//ff:what IsUnderBudget: 요청 금액이 잔여 예산 이하인지 판정
package approve

// IsUnderBudget checks if the requested amount is within remaining budget.
func IsUnderBudget(claim any, ground any, backing any) (bool, any) {
	req := claim.(*ApprovalRequest)
	ctx := ground.(*ApprovalContext)
	return req.Amount <= ctx.Budget.Remaining, nil
}
