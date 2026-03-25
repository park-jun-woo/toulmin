//ff:func feature=approve type=rule control=sequence
//ff:what HasApprovalRole: backing(ApproverBacking)으로 지정된 결재 역할을 가졌는지 판정
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// HasApprovalRole checks if the approver has the role specified by backing.
func HasApprovalRole(ctx toulmin.Context, backing toulmin.Backing) (bool, any) {
	approverRole, _ := ctx.Get("approverRole")
	ab := backing.(*ApproverBacking)
	return approverRole.(string) == ab.Role, nil
}
