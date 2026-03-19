//ff:type feature=policy type=model
//ff:what RequestContext: 정책 판정에 필요한 요청 컨텍스트 (런타임 데이터만)
package policy

// RequestContext holds per-request facts for policy evaluation.
// User is any — the framework does not impose a concrete User type.
// Field access is done via backing (extraction functions).
type RequestContext struct {
	User            any
	ClientIP        string
	ResourceOwnerID string
	Headers         map[string]string
	RateLimiter     RateLimiter
	Metadata        map[string]any
}
