//ff:type feature=route type=model
//ff:what User: 라우트 판정에 사용되는 사용자 정보
package route

// User represents user identity for route evaluation.
type User struct {
	ID    string
	Role  string
	Email string
}
