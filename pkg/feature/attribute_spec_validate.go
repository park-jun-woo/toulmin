//ff:func feature=feature type=model control=sequence
//ff:what AttributeSpec.Validate: spec 필드 유효성 검증
package feature

import "fmt"

// Validate checks that AttributeSpec fields are valid.
func (b *AttributeSpec) Validate() error {
	if b.Key == "" {
		return fmt.Errorf("AttributeSpec: Key must not be empty")
	}
	return nil
}
