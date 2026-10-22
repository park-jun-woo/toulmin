//ff:func feature=tangl type=util control=sequence
//ff:what readIDWidth — read Meta ID-WIDTH into the Edge Id width, enforcing the 10 reserved bits
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/core"
)

// readIDWidth reads the 1-word ID-WIDTH payload `EE VV II 0000000000` (GEUL t002
// OVERVIEW 정의): EE is the Edge Id width this profile uses. VV and II are
// accepted as declared; the 10 low bits are reserved and must be 0.
func readIDWidth(h *headInfo, r core.Record) error {
	if len(r.Payload) != 1 || r.Payload[0]&0x03FF != 0 {
		return errors.New("tangeul/stream: ID-WIDTH must be one word with reserved bits 0")
	}
	h.idWidth = core.WidthWords(uint8(r.Payload[0] >> 14))
	return nil
}
