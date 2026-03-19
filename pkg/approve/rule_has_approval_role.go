//ff:func feature=approve type=rule control=sequence
//ff:what HasApprovalRole: backing(string)으로 지정된 결재 역할을 가졌는지 판정
package approve

// HasApprovalRole checks if the approver has the role specified by backing (string).
func HasApprovalRole(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ApprovalContext)
	role := backing.(string)
	return ctx.Approver.Role == role, nil
}
