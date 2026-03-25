//ff:func feature=moderate type=rule control=sequence
//ff:what IsTrustedUser: spec(TrustScoreSpec)으로 지정된 최소 신뢰 점수 이상인지 판정
package moderate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// IsTrustedUser checks if the author's trust score meets the minimum.
// spec is *TrustScoreSpec.
func IsTrustedUser(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
	author, _ := ctx.Get("author")
	tb := specs[0].(*TrustScoreSpec)
	return author.(*Author).TrustScore >= tb.MinScore, nil
}
