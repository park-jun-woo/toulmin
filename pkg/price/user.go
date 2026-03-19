//ff:type feature=price type=model
//ff:what User: 사용자 정보
package price

// User represents user identity for price evaluation.
type User struct {
	ID         string
	Membership string // "none", "basic", "gold", "vip"
}
