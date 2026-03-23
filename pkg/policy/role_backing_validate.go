//ff:func feature=policy type=model control=sequence
//ff:what RoleBacking.Validate: 필수 필드 검증
package policy

import "fmt"

// Validate checks that Role is non-empty.
func (b *RoleBacking) Validate() error {
	if b.Role == "" {
		return fmt.Errorf("RoleBacking: Role is required")
	}
	return nil
}
