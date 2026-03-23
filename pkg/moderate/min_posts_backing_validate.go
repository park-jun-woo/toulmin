//ff:func feature=moderate type=model control=sequence
//ff:what MinPostsBacking.Validate: backing 필드 유효성 검증
package moderate

import "fmt"

// Validate checks that MinPostsBacking fields are valid.
func (b *MinPostsBacking) Validate() error {
	if b.MinPosts < 0 {
		return fmt.Errorf("MinPostsBacking: MinPosts must be non-negative, got %d", b.MinPosts)
	}
	return nil
}
