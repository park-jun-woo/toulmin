//ff:func feature=policy type=model control=sequence
//ff:what RoleSpec.Validate: 필수 필드 검증
package policy

import "fmt"

// Validate checks that Role is non-empty.
func (b *RoleSpec) Validate() error {
	if b.Role == "" {
		return fmt.Errorf("RoleSpec: Role is required")
	}
	return nil
}
