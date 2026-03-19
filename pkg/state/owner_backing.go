//ff:type feature=state type=model
//ff:what OwnerBacking: IsOwner rule의 backing 타입 (ID 추출 함수 쌍)
package state

// OwnerBacking carries ID extraction functions for ownership checks.
type OwnerBacking struct {
	OwnerIDFunc func(any) string
	UserIDFunc  func(any) string
}
