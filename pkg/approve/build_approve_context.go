//ff:func feature=approve type=adapter control=sequence
//ff:what buildApproveContext: ApprovalRequest + ApprovalContext → toulmin.Context 변환
package approve

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildApproveContext converts an ApprovalRequest and ApprovalContext into a toulmin.Context.
func buildApproveContext(req *ApprovalRequest, ac *ApprovalContext) toulmin.Context {
	ctx := toulmin.NewContext()
	ctx.Set("amount", req.Amount)
	ctx.Set("category", req.Category)
	ctx.Set("requesterID", req.RequesterID)
	ctx.Set("description", req.Description)
	ctx.Set("requestMetadata", req.Metadata)
	ctx.Set("approver", ac.Approver)
	ctx.Set("approverID", ac.ApproverID)
	ctx.Set("approverRole", ac.ApproverRole)
	ctx.Set("approverLevel", ac.ApproverLevel)
	ctx.Set("budget", ac.Budget)
	ctx.Set("orgTree", ac.OrgTree)
	ctx.Set("metadata", ac.Metadata)
	return ctx
}
