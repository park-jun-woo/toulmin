//ff:func feature=tangl type=util control=sequence
//ff:what expectSpecVersion — require Meta SPEC-VERSION to name the spec this build reads
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// expectSpecVersion requires r, Meta SPEC-VERSION, to carry toulmin.SpecVersion
// as its bare unsigned integer.
func expectSpecVersion(r core.Record) error {
	v, err := core.DecodeSpecVersion(r)
	if err != nil {
		return err
	}
	if v != toulmin.SpecVersion {
		return fmt.Errorf("tangeul/stream: SPEC-VERSION is 0x%04x, this build reads 0x%04x", v, toulmin.SpecVersion)
	}
	return nil
}
