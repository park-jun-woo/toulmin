//ff:func feature=moderate type=model control=sequence
//ff:what ClassifierSpec.Validate: spec 필드 유효성 검증
package moderate

import "fmt"

// Validate checks that ClassifierSpec fields are valid.
func (b *ClassifierSpec) Validate() error {
	if b.Classifier == nil {
		return fmt.Errorf("ClassifierSpec: Classifier must not be nil")
	}
	return nil
}
