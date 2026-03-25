//ff:func feature=price type=model control=sequence
//ff:what MemberSpec.Validate: spec 필드 유효성 검증
package price

import "fmt"

// Validate checks that MemberSpec fields are valid.
func (b *MemberSpec) Validate() error {
	if b.Level == "" {
		return fmt.Errorf("MemberSpec: Level must not be empty")
	}
	return nil
}
