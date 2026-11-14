//ff:type feature=approve type=interface
//ff:what OrgTree: 조직 구조 추상화 인터페이스
package approve

// OrgTree abstracts organizational hierarchy queries.
type OrgTree interface {
	IsDirectManager(approverID, requesterID string) bool
	Level(userID string) int
}
