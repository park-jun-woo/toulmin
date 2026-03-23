//ff:func feature=feature type=model control=sequence
//ff:what PercentageBacking.Validate: backing 필드 유효성 검증
package feature

import "fmt"

// Validate checks that PercentageBacking fields are valid.
func (b *PercentageBacking) Validate() error {
	if b.Percentage < 0 || b.Percentage > 1 {
		return fmt.Errorf("PercentageBacking: Percentage must be between 0 and 1, got %f", b.Percentage)
	}
	return nil
}
