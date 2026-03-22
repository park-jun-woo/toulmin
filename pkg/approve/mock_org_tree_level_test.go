//ff:func feature=approve type=model control=sequence
//ff:what mockOrgTree.Level — returns level for user (stub)
package approve

func (m *mockOrgTree) Level(userID string) int { return 0 }
