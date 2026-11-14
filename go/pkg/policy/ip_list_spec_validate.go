//ff:func feature=policy type=model control=sequence
//ff:what IPListSpec.Validate: 필수 필드 검증
package policy

import "fmt"

// Validate checks that Purpose is non-empty.
func (b *IPListSpec) Validate() error {
	if b.Purpose == "" {
		return fmt.Errorf("IPListSpec: Purpose is required")
	}
	return nil
}
