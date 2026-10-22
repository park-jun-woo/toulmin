//ff:func feature=feature type=model control=sequence
//ff:what PercentageSpec.Validate: spec 필드 유효성 검증
package feature

import "fmt"

// Validate checks that PercentageSpec fields are valid.
func (b *PercentageSpec) Validate() error {
	if b.Percentage < 0 || b.Percentage > 1 {
		return fmt.Errorf("PercentageSpec: Percentage must be between 0 and 1, got %f", b.Percentage)
	}
	return nil
}
