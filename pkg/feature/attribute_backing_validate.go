//ff:func feature=feature type=model control=sequence
//ff:what AttributeBacking.Validate: backing 필드 유효성 검증
package feature

import "fmt"

// Validate checks that AttributeBacking fields are valid.
func (b *AttributeBacking) Validate() error {
	if b.Key == "" {
		return fmt.Errorf("AttributeBacking: Key must not be empty")
	}
	return nil
}
