//ff:type feature=moderate type=model
//ff:what Channel: 콘텐츠 게시 채널 정보
package moderate

// Channel represents the channel where content is posted.
type Channel struct {
	ID       string
	Type     string
	AgeGated bool
}
