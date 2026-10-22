//ff:type feature=state type=model
//ff:what TransitionContext: 전이 판정에 필요한 런타임 컨텍스트
package state

// TransitionContext holds per-request facts for transition evaluation.
type TransitionContext struct {
	CurrentState    string
	User            any
	Resource        any
	Metadata        map[string]any
	UserID          string
	ResourceOwnerID string
}
