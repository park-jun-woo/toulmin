//ff:func feature=feature type=model control=sequence
//ff:what RegionBacking.Validate: backing 필드 유효성 검증
package feature

import "fmt"

// Validate checks that RegionBacking fields are valid.
func (b *RegionBacking) Validate() error {
	if b.Region == "" {
		return fmt.Errorf("RegionBacking: Region must not be empty")
	}
	return nil
}
