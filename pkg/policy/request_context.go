//ff:type feature=policy type=model
//ff:what RequestContext: 정책 판정에 필요한 요청 컨텍스트 (런타임 데이터만)
package policy

// RequestContext holds per-request facts for policy evaluation.
// Judgment criteria (IP lists, role names, thresholds) belong in backing, not here.
type RequestContext struct {
	User            *User
	ClientIP        string
	ResourceOwnerID string
	Headers         map[string]string
	RateLimiter     RateLimiter
	Metadata        map[string]any
}
