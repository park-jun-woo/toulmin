//ff:func feature=price type=model control=sequence
//ff:what BulkOrderBacking.Validate: backing 필드 유효성 검증
package price

import "fmt"

// Validate checks that BulkOrderBacking fields are valid.
func (b *BulkOrderBacking) Validate() error {
	if b.MinQuantity <= 0 {
		return fmt.Errorf("BulkOrderBacking: MinQuantity must be positive, got %d", b.MinQuantity)
	}
	return nil
}
