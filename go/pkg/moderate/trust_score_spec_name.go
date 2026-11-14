//ff:func feature=moderate type=model control=sequence
//ff:what TrustScoreSpec.SpecName: spec 타입 식별자 반환
package moderate

// SpecName returns the type identifier for TrustScoreSpec.
func (b *TrustScoreSpec) SpecName() string {
	return "TrustScoreSpec"
}
