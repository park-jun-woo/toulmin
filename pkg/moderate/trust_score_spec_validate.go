//ff:func feature=moderate type=model control=sequence
//ff:what TrustScoreSpec.Validate: spec 필드 유효성 검증
package moderate

import "fmt"

// Validate checks that TrustScoreSpec fields are valid.
func (b *TrustScoreSpec) Validate() error {
	if b.MinScore < 0 || b.MinScore > 1 {
		return fmt.Errorf("TrustScoreSpec: MinScore must be between 0 and 1, got %f", b.MinScore)
	}
	return nil
}
