//ff:type feature=policy type=model
//ff:what User: 정책 판정에 사용되는 사용자 정보
package policy

// User represents user identity for policy evaluation.
type User struct {
	ID    string
	Role  string
	Email string
}
