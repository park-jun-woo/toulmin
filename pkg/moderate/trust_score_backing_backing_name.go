//ff:func feature=moderate type=model control=sequence
//ff:what TrustScoreBacking.BackingName: backing 타입 식별자 반환
package moderate

// BackingName returns the type identifier for TrustScoreBacking.
func (b *TrustScoreBacking) BackingName() string {
	return "TrustScoreBacking"
}
