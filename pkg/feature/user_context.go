//ff:type feature=feature type=model
//ff:what UserContext: 피처 판정에 필요한 사용자 컨텍스트
package feature

// UserContext holds per-request facts for feature evaluation.
// User is any — the framework does not impose a concrete User type.
type UserContext struct {
	User       any
	ID         string
	Region     string
	Attributes map[string]any
}
