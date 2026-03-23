//ff:func feature=moderate type=model control=sequence
//ff:what TrustScoreBacking.Validate: backing 필드 유효성 검증
package moderate

import "fmt"

// Validate checks that TrustScoreBacking fields are valid.
func (b *TrustScoreBacking) Validate() error {
	if b.MinScore < 0 || b.MinScore > 1 {
		return fmt.Errorf("TrustScoreBacking: MinScore must be between 0 and 1, got %f", b.MinScore)
	}
	return nil
}
