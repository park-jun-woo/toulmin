//ff:func feature=feature type=model control=sequence
//ff:what RegionSpec.Validate: spec 필드 유효성 검증
package feature

import "fmt"

// Validate checks that RegionSpec fields are valid.
func (b *RegionSpec) Validate() error {
	if b.Region == "" {
		return fmt.Errorf("RegionSpec: Region must not be empty")
	}
	return nil
}
