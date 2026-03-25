//ff:func feature=moderate type=rule control=sequence
//ff:what IsTrustedUser: backing(TrustScoreBacking)으로 지정된 최소 신뢰 점수 이상인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsTrustedUser checks if the author's trust score meets the minimum.
// backing is *TrustScoreBacking.
func IsTrustedUser(claim any, ground any, backing toulmin.Backing) (bool, any) {
	ctx := ground.(*ContentContext)
	tb := backing.(*TrustScoreBacking)
	return ctx.Author.TrustScore >= tb.MinScore, nil
}
