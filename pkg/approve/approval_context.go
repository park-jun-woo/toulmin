//ff:type feature=approve type=model
//ff:what ApprovalContext: 승인 판정에 필요한 런타임 컨텍스트
package approve

// ApprovalContext holds per-step facts for approval evaluation.
// Approver is any — the framework does not impose a concrete Approver type.
type ApprovalContext struct {
	Approver any
	Budget   *Budget
	OrgTree  OrgTree
	Metadata map[string]any
}
