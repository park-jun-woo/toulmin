//ff:func feature=price type=model control=sequence
//ff:what DiscountBacking.Validate: backing 필드 유효성 검증
package price

import "fmt"

// Validate checks that DiscountBacking fields are valid.
func (b *DiscountBacking) Validate() error {
	if b.Rate < 0 || b.Rate > 1 {
		return fmt.Errorf("DiscountBacking: Rate must be between 0 and 1, got %f", b.Rate)
	}
	if b.Fixed < 0 {
		return fmt.Errorf("DiscountBacking: Fixed must be non-negative, got %f", b.Fixed)
	}
	return nil
}
