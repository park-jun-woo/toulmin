//ff:type feature=moderate type=model
//ff:what TrustScoreBacking: IsTrustedUser rule의 backing 타입
package moderate

// TrustScoreBacking carries minimum trust score criteria.
type TrustScoreBacking struct {
	MinScore float64 // minimum trust score threshold
}
