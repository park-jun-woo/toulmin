//ff:func feature=approve type=model control=sequence
//ff:what mockOrgTree.IsDirectManager — checks if approver is direct manager of requester
package approve

func (m *mockOrgTree) IsDirectManager(approverID, requesterID string) bool {
	return m.managers[requesterID] == approverID
}
