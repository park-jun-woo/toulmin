//ff:type feature=moderate type=model
//ff:what Author: 콘텐츠 작성자 정보
package moderate

// Author represents the content author.
type Author struct {
	ID         string
	Verified   bool
	PostCount  int
	TrustScore float64
}
