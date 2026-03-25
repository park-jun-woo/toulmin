//ff:func feature=price type=model control=sequence
//ff:what BulkOrderSpec.Validate: spec 필드 유효성 검증
package price

import "fmt"

// Validate checks that BulkOrderSpec fields are valid.
func (b *BulkOrderSpec) Validate() error {
	if b.MinQuantity <= 0 {
		return fmt.Errorf("BulkOrderSpec: MinQuantity must be positive, got %d", b.MinQuantity)
	}
	return nil
}
