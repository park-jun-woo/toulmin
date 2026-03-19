//ff:func feature=moderate type=rule control=sequence
//ff:what IsTrustedUser: backing(float64)으로 지정된 최소 신뢰 점수 이상인지 판정
package moderate

// IsTrustedUser checks if the author's trust score meets the minimum.
// backing is float64 (minimum trust score).
func IsTrustedUser(claim any, ground any, backing any) (bool, any) {
	ctx := ground.(*ContentContext)
	minScore := backing.(float64)
	return ctx.Author.TrustScore >= minScore, nil
}
