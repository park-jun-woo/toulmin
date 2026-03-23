//ff:func feature=price type=model control=sequence
//ff:what MemberBacking.Validate: backing 필드 유효성 검증
package price

import "fmt"

// Validate checks that MemberBacking fields are valid.
func (b *MemberBacking) Validate() error {
	if b.Level == "" {
		return fmt.Errorf("MemberBacking: Level must not be empty")
	}
	return nil
}
