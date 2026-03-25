//ff:type feature=moderate type=model
//ff:what TrustScoreSpec: IsTrustedUser rule의 spec 타입
package moderate

// TrustScoreSpec carries minimum trust score criteria.
type TrustScoreSpec struct {
	MinScore float64 // minimum trust score threshold
}
