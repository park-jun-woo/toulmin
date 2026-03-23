//ff:func feature=policy type=model control=sequence
//ff:what IPListBacking.Validate: 필수 필드 검증
package policy

import "fmt"

// Validate checks that Purpose is non-empty.
func (b *IPListBacking) Validate() error {
	if b.Purpose == "" {
		return fmt.Errorf("IPListBacking: Purpose is required")
	}
	return nil
}
