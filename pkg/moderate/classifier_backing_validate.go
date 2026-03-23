//ff:func feature=moderate type=model control=sequence
//ff:what ClassifierBacking.Validate: backing 필드 유효성 검증
package moderate

import "fmt"

// Validate checks that ClassifierBacking fields are valid.
func (b *ClassifierBacking) Validate() error {
	if b.Classifier == nil {
		return fmt.Errorf("ClassifierBacking: Classifier must not be nil")
	}
	return nil
}
