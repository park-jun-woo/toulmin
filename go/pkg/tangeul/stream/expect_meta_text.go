//ff:func feature=tangl type=util control=sequence
//ff:what expectMetaText — require a Meta record's text Value to equal want
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/core"
)

// expectMetaText requires that r's single text Value equals want.
func expectMetaText(r core.Record, evenCodec uint8, want string) error {
	s, err := metaText(r, evenCodec)
	if err != nil {
		return err
	}
	if s != want {
		return fmt.Errorf("tangeul/stream: Meta 0x%02X is %q, this build reads %q", r.Codec, s, want)
	}
	return nil
}
