//ff:func feature=price type=model control=sequence
//ff:what DiscountSpec.Validate: spec 필드 유효성 검증
package price

import "fmt"

// Validate checks that DiscountSpec fields are valid.
func (b *DiscountSpec) Validate() error {
	if b.Rate < 0 || b.Rate > 1 {
		return fmt.Errorf("DiscountSpec: Rate must be between 0 and 1, got %f", b.Rate)
	}
	if b.Fixed < 0 {
		return fmt.Errorf("DiscountSpec: Fixed must be non-negative, got %f", b.Fixed)
	}
	return nil
}
