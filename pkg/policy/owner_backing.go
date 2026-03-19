//ff:type feature=policy type=model
//ff:what OwnerBacking: IsOwner rule의 backing 타입 (ID 추출 함수 쌍)
package policy

// OwnerBacking carries ID extraction functions for ownership checks.
type OwnerBacking struct {
	UserIDFunc    func(any) string
	ResourceIDFunc func(any) string
}
