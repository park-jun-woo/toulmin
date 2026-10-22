//ff:type feature=moderate type=model
//ff:what ContentContext: 모더레이션 판정에 필요한 런타임 컨텍스트
package moderate

// ContentContext holds per-request facts for moderation evaluation.
type ContentContext struct {
	Author   *Author
	Channel  *Channel
	Metadata map[string]any
}
