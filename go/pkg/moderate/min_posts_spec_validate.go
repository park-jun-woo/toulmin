//ff:func feature=moderate type=model control=sequence
//ff:what MinPostsSpec.Validate: spec 필드 유효성 검증
package moderate

import "fmt"

// Validate checks that MinPostsSpec fields are valid.
func (b *MinPostsSpec) Validate() error {
	if b.MinPosts < 0 {
		return fmt.Errorf("MinPostsSpec: MinPosts must be non-negative, got %d", b.MinPosts)
	}
	return nil
}
