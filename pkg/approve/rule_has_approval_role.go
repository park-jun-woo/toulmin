//ff:func feature=approve type=rule control=sequence
//ff:what HasApprovalRole: spec(ApproverSpec)으로 지정된 결재 역할을 가졌는지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasApprovalRole checks if the approver has the role specified by spec.
func HasApprovalRole(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	approverRole, _ := ctx.Get("approverRole")
	if len(specs) == 0 {
		return false, nil
	}
	ab := specs[0].(*ApproverSpec)
	return approverRole.(string) == ab.Role, nil
}
